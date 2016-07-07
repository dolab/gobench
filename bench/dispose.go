package bench

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type DisposeWorkerIn struct {
	Bucket string
}

type DisposeWorkerOut struct {
	Ok     bool
	Cost   time.Duration
	Bucket string
	Err    string
}

func (out *DisposeWorkerOut) String() string {
	return fmt.Sprintf(`{"ok":%t,"bucket":"%s","err":"%s"} %v`, out.Ok, out.Bucket, out.Err, out.Cost)
}

// for bucket creation
func DisposeWorker(n int, cli *s3.S3, jobs <-chan *DisposeWorkerIn, results chan<- bool) {
	for in := range jobs {
		out := DeleteBucket(cli, in.Bucket)
		if out.Ok {
			blog.Print("[#", n, "] ", out.String())
		} else {
			blog.Error("[#", n, "] ", out.String())
		}

		results <- true
	}
}

func DeleteBucket(cli *s3.S3, bucket string) *DisposeWorkerOut {
	s3payload := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}

	// download
	t := time.Now()
	_, s3err := cli.DeleteBucket(s3payload)

	result := &DisposeWorkerOut{
		Ok:     (s3err == nil),
		Cost:   time.Since(t),
		Bucket: bucket,
	}
	if !result.Ok {
		result.Err = s3err.Error()
	}

	return result
}
