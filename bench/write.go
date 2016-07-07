package bench

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/satori/go.uuid"
)

type WriteWorkerIn struct {
	Bucket string
	Key    string
	File   string
	Size   int
}

type WriteWorkerOut struct {
	Ok   bool
	Cost time.Duration
	Key  string
	Err  string
}

func (out *WriteWorkerOut) String() string {
	return fmt.Sprintf(`{"ok":%t,"key":"%s","err":"%s"} %v`, out.Ok, out.Key, out.Err, out.Cost)
}

func WriteWorker(n int, cli *s3.S3, hashCheck bool, jobs <-chan *WriteWorkerIn, results chan<- bool) {
	for in := range jobs {
		s3file := in.File
		if s3file == "" {
			out := &WriteWorkerOut{
				Ok:  false,
				Key: "",
			}

			// create a temporary file for uploading
			tmpfile, tmperr := ioutil.TempFile("", "upload")
			if tmperr != nil {
				out.Err = tmperr.Error()

				blog.Error("[#", n, "] ", out.String())
				return
			}

			blockSize := 32
			blockContent := uuid.NewV4().String()[:blockSize]
			for i := 0; i < in.Size/blockSize; i++ {
				n, err := tmpfile.WriteString(blockContent)
				if n != blockSize || err != nil {
					out.Err = fmt.Sprintf("tmpfile.WriteString(), expected %d, but got %d with %v", blockSize, n, err)

					blog.Error("[#", n, "] ", out.String())
					return
				}
			}

			s3file = tmpfile.Name()
			defer os.Remove(s3file)
		}

		out := PutObject(cli, in.Bucket, in.Key, s3file, hashCheck)
		if out.Ok {
			blog.Print("[#", n, "] ", out.String())
		} else {
			blog.Error("[#", n, "] ", out.String())
		}

		results <- true
	}
}

func PutObject(cli *s3.S3, bucket, key, file string, hashCheck bool) *WriteWorkerOut {
	result := &WriteWorkerOut{
		Ok:  false,
		Key: bucket + ":" + key,
	}

	freader, err := os.Open(file)
	if err != nil {
		result.Err = err.Error()

		return result
	}
	defer freader.Close()

	// resolve file mime by file extension
	// s3mime := mime.TypeByExtension(path.Ext(file))
	// if s3mime == "" {
	// 	s3mime = "application/octet-stream"
	// }
	s3mime := "application/octet-stream"

	// calc hash
	md5hash := md5.New()
	if hashCheck {
		_, err := io.Copy(md5hash, freader)
		if err != nil {
			result.Err = err.Error()

			return result
		}

		freader.Seek(0, 0)
	}

	s3payload := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(s3mime),
		Body:        freader,
	}

	// upload
	t := time.Now()
	s3result, s3err := cli.PutObject(s3payload)

	result.Ok = (s3err == nil)
	result.Cost = time.Since(t)
	if !result.Ok {
		result.Err = s3err.Error()

		return result
	}

	if hashCheck {
		hash := `"` + hex.EncodeToString(md5hash.Sum(nil)) + `"`
		if strings.Compare(*s3result.ETag, hash) != 0 {
			result.Err = "[ETAG] write expected " + hash + ", but got " + *s3result.ETag
		}
	}

	return result
}
