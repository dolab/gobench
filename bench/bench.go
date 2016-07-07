package bench

import (
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dolab/logger"
)

var (
	blog *logger.Logger
)

func StartWorkflow(s3service *s3.S3, works []*WorkConfig, logger *logger.Logger) {
	if len(works) == 0 {
		logger.Fatal("Workflow CAN NOT be empty!")

		return
	}

	blog = logger

	for _, work := range works {
		if !work.Enabled {
			blog.Debugf("Skiped work %v", work)
			continue
		}

		switch work.Stage {
		case StageInit:
			var (
				inchan  = make(chan *InitWorkerIn, work.Concurrent)
				outchan = make(chan bool, work.BucketNumber())
			)

			// start workers
			for i := 0; i < work.Concurrent; i++ {
				go InitWorker(i+1, s3service, inchan, outchan)
			}

			// emit jobs
			var (
				ticker  = time.NewTicker(time.Second / time.Duration(work.Concurrent))
				counter = work.BucketStart
				breaked = false
			)
			for {
				select {
				case <-ticker.C:
					inchan <- &InitWorkerIn{
						Bucket: work.BucketPrefix + strconv.Itoa(counter),
					}

					if counter == work.BucketEnd {
						ticker.Stop()

						breaked = true
					}

					counter += 1
				}

				if breaked {
					break
				}
			}

			close(inchan)

			// waiting for all workers finished
			for n := 0; n < work.BucketNumber(); n++ {
				<-outchan
			}

		case StageWrite:
			if work.Filesize < 0 {
				work.Filesize = 0
			}

			var (
				inchan  = make(chan *WriteWorkerIn, work.Concurrent)
				outchan = make(chan bool, work.ObjectTotals())
			)

			// start workers
			for i := 0; i < work.Concurrent; i++ {
				go WriteWorker(i+1, s3service, work.HashCheck, inchan, outchan)
			}

			// emit jobs
			delta := time.Second / time.Duration(work.Concurrent)
			for j := 1; j <= work.BucketNumber(); j++ {
				// reset
				var (
					ticker  = time.NewTicker(delta)
					counter = work.ObjectStart
					breaked = false
				)

				for {
					select {
					case <-ticker.C:
						inchan <- &WriteWorkerIn{
							Bucket: work.BucketPrefix + strconv.Itoa(j),
							Key:    work.ObjectPrefix + strconv.Itoa(counter),
							Size:   work.Filesize,
						}

						if counter == work.ObjectEnd {
							ticker.Stop()

							breaked = true
						}

						counter += 1
					}

					if breaked {
						break
					}
				}
			}

			close(inchan)

			// waiting for all workers finished
			for j := 0; j < work.ObjectTotals(); j++ {
				<-outchan
			}

		case StageRead:
			var (
				inchan  = make(chan *ReadWorkerIn, work.Concurrent)
				outchan = make(chan bool, work.ObjectTotals())
			)

			// start workers
			for i := 0; i < work.Concurrent; i++ {
				go ReadWorker(i+1, s3service, work.HashCheck, inchan, outchan)
			}

			// emit jobs
			delta := time.Second / time.Duration(work.Concurrent)
			for j := 1; j <= work.BucketNumber(); j++ {
				// reset
				var (
					ticker  = time.NewTicker(delta)
					counter = work.ObjectStart
					breaked = false
				)

				for {
					select {
					case <-ticker.C:
						inchan <- &ReadWorkerIn{
							Bucket: work.BucketPrefix + strconv.Itoa(j),
							Key:    work.ObjectPrefix + strconv.Itoa(counter),
						}

						if counter == work.ObjectEnd {
							ticker.Stop()

							breaked = true
						}

						counter += 1
					}

					if breaked {
						break
					}
				}
			}

			close(inchan)

			// waiting for all workers finished
			for n := 0; n < work.ObjectTotals(); n++ {
				<-outchan
			}

		case StageClean:
			var (
				inchan  = make(chan *CleanWorkerIn, work.Concurrent)
				outchan = make(chan bool, work.ObjectTotals())
			)

			// start workers
			for i := 0; i < work.Concurrent; i++ {
				go CleanWorker(i+1, s3service, inchan, outchan)
			}

			// emit jobs
			delta := time.Second / time.Duration(work.Concurrent)
			for j := 1; j <= work.BucketNumber(); j++ {
				var (
					ticker  = time.NewTicker(delta)
					counter = work.ObjectStart
					breaked = false
				)

				for {
					select {
					case <-ticker.C:
						inchan <- &CleanWorkerIn{
							Bucket: work.BucketPrefix + strconv.Itoa(j),
							Key:    work.ObjectPrefix + strconv.Itoa(counter),
						}

						if counter == work.ObjectEnd {
							ticker.Stop()

							breaked = true
						}

						counter += 1
					}

					if breaked {
						break
					}
				}
			}

			close(inchan)

			// waiting for all workers finished
			for n := 0; n < work.ObjectTotals(); n++ {
				<-outchan
			}

		case StageDispose:
			var (
				inchan  = make(chan *DisposeWorkerIn, work.Concurrent)
				outchan = make(chan bool, work.BucketNumber())
			)

			// start workers
			for i := 0; i < work.Concurrent; i++ {
				go DisposeWorker(i+1, s3service, inchan, outchan)
			}

			// emit jobs
			var (
				ticker  = time.NewTicker(time.Second / time.Duration(work.Concurrent))
				counter = work.BucketStart
				breaked = false
			)
			for {
				select {
				case <-ticker.C:
					inchan <- &DisposeWorkerIn{
						Bucket: work.BucketPrefix + strconv.Itoa(counter),
					}

					if counter == work.BucketEnd {
						ticker.Stop()

						breaked = true
					}

					counter += 1
				}

				if breaked {
					break
				}
			}

			close(inchan)

			// waiting for all workers finished
			for j := 0; j < work.BucketNumber(); j++ {
				<-outchan
			}

		default:
			log.Fatal("Unsupported stage mode: ", work.Stage)
		}
	}
}
