package main

import (
	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dolab/bench"
	"github.com/dolab/logger"
)

var (
	s3config string
	s3log    string
	s3debug  bool
)

func main() {
	flag.StringVar(&s3config, "s3.config", "", "-s3.config=/path/to/config.json")
	flag.StringVar(&s3log, "s3.log", "", "-s3.log=/path/to/gobench.log")
	flag.BoolVar(&s3debug, "s3.debug", false, "-s3.debug=true|false")
	flag.Parse()

	if s3config == "" {
		flag.PrintDefaults()
		return
	}

	config, err := bench.NewConfig(s3config)
	if err != nil {
		panic(err.Error())
	}

	s3provider := credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     config.AccessKeyId,
			SecretAccessKey: config.AccessKeySecret,
		},
	}

	s3credential := credentials.NewCredentials(&s3provider)

	s3session := session.New()
	s3session.Config.WithEndpoint(config.Host)
	s3session.Config.WithRegion(config.Region)
	s3session.Config.WithCredentials(s3credential)
	s3session.Config.WithMaxRetries(1)
	s3session.Config.WithS3ForcePathStyle(true)
	if s3debug == true {
		s3session.Config.WithLogLevel(aws.LogDebug)
	}

	s3service := s3.New(s3session, nil)

	blog, _ := logger.New("stderr")
	blog.SetColor(false)

	bench.StartWorkflow(s3service, config.Workflow, blog)
}
