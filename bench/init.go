package bench

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type InitWorkerIn struct {
	Bucket string
}

type InitWorkerOut struct {
	Ok     bool
	Cost   time.Duration
	Bucket string
	Err    string
}

func (out *InitWorkerOut) String() string {
	return fmt.Sprintf(`{"ok":%t,"bucket":"%s","err":"%s"} %v`, out.Ok, out.Bucket, out.Err, out.Cost)
}

// for bucket creation
func InitWorker(n int, cli *s3.S3, jobs <-chan *InitWorkerIn, result chan<- bool) {
	for in := range jobs {
		out := CreateBucket(cli, in.Bucket)
		if out.Ok {
			blog.Print("[#", n, "] ", out.String())
		} else {
			blog.Error("[#", n, "] ", out.String())
		}

		result <- true
	}
}

func CreateBucket(cli *s3.S3, bucket string) *InitWorkerOut {
	s3payload := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	// download
	t := time.Now()
	_, s3err := cli.CreateBucket(s3payload)

	result := &InitWorkerOut{
		Ok:     (s3err == nil),
		Cost:   time.Since(t),
		Bucket: bucket,
	}
	if !result.Ok {
		result.Err = s3err.Error()
	}

	return result
}
