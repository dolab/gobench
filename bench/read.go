package bench

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ReadWorkerIn struct {
	Bucket string
	Key    string
}

type ReadWorkerOut struct {
	Ok   bool
	Cost time.Duration
	Key  string
	Err  string
}

func (out *ReadWorkerOut) String() string {
	return fmt.Sprintf(`{"ok":%t,"key":"%s","err":"%s"} %v`, out.Ok, out.Key, out.Err, out.Cost)
}

func ReadWorker(n int, cli *s3.S3, hashCheck bool, jobs <-chan *ReadWorkerIn, results chan<- bool) {
	for in := range jobs {
		out := GetObject(cli, in.Bucket, in.Key, hashCheck)
		if out.Ok {
			blog.Print("[#", n, "] ", out.String())
		} else {
			blog.Error("[#", n, "] ", out.String())
		}

		results <- true
	}
}

func GetObject(cli *s3.S3, bucket, key string, hashCheck bool) *ReadWorkerOut {
	s3payload := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	// download
	t := time.Now()
	s3result, s3err := cli.GetObject(s3payload)

	result := &ReadWorkerOut{
		Ok:   (s3err == nil),
		Cost: time.Since(t),
		Key:  bucket + ":" + key,
	}
	if !result.Ok {
		result.Err = s3err.Error()

		return result
	}
	defer s3result.Body.Close()

	if hashCheck {
		md5hash := md5.New()
		io.Copy(md5hash, s3result.Body)

		hash := `"` + hex.EncodeToString(md5hash.Sum(nil)) + `"`
		if *s3result.ETag != hash {
			result.Err = "[ETAG] read expected " + *s3result.ETag + ", but got " + hash
		}
	}

	return result
}
