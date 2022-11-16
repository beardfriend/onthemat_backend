package aws

import (
	"context"
	"io"
	"log"

	localCnf "onthemat/internal/app/config"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type S3 interface {
	SetConfig() error
	Upload(key string, file io.ReadSeeker) *manager.UploadOutput
}
type s3 struct {
	AwsS3Region  string
	AwsAccessKey string
	AwsSecretKey string
	BucketName   string
	client       *awss3.Client
}

func NewS3(localCnf *localCnf.Config) S3 {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := awss3.NewFromConfig(cfg)

	return &s3{
		client:       client,
		AwsS3Region:  localCnf.AWSS3.Region,
		BucketName:   localCnf.AWSS3.BucketName,
		AwsAccessKey: localCnf.AWS.AceessKey,
		AwsSecretKey: localCnf.AWS.SecretKey,
	}
}

func (s *s3) SetConfig() error {
	creds := credentials.NewStaticCredentialsProvider(s.AwsAccessKey, s.AwsSecretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds),
		config.WithRegion(s.AwsS3Region),
	)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	s.client = awss3.NewFromConfig(cfg)
	return nil
}

func (s *s3) Upload(key string, file io.ReadSeeker) *manager.UploadOutput {
	uploader := manager.NewUploader(s.client)
	result, err := uploader.Upload(context.TODO(), &awss3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return result
}
