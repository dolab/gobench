package bench

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type CleanWorkerIn struct {
	Bucket string
	Key    string
}

type CleanWorkerOut struct {
	Ok   bool
	Cost time.Duration
	Key  string
	Err  string
}

func (out *CleanWorkerOut) String() string {
	return fmt.Sprintf(`{"ok":%t,"key":"%s","err":"%s"} %v`, out.Ok, out.Key, out.Err, out.Cost)
}

// for bucket creation
func CleanWorker(n int, cli *s3.S3, jobs <-chan *CleanWorkerIn, results chan<- bool) {
	for in := range jobs {
		out := DeleteObject(cli, in.Bucket, in.Key)
		if out.Ok {
			blog.Print("[#", n, "] ", out.String())
		} else {
			blog.Error("[#", n, "] ", out.String())
		}

		results <- true
	}
}

func DeleteObject(cli *s3.S3, bucket, key string) *CleanWorkerOut {
	s3payload := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// download
	t := time.Now()
	_, s3err := cli.DeleteObject(s3payload)

	result := &CleanWorkerOut{
		Ok:   (s3err == nil),
		Cost: time.Since(t),
		Key:  bucket + ":" + key,
	}
	if !result.Ok {
		result.Err = s3err.Error()
	}

	return result
}
