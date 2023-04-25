package amazonS3

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AmazonS3 struct {
	Region          string `env:""`
	AccessKeyID     string `env:""`
	SecretAccessKey string `env:""`
	SessionToken    string `env:""`
	BucketName      string `env:""`

	cli *s3.S3
}

func NewAmazonS3(regin, accessKeyID, secretAccessKey, sessionToken, bucketName string) *AmazonS3 {
	s3 := &AmazonS3{
		Region:          regin,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		SessionToken:    sessionToken,
		BucketName:      bucketName,
	}
	if err := s3.Init(); err != nil {
		panic(err)
	}
	return s3
}

func (s *AmazonS3) Init() error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(s.Region),
		Credentials: credentials.NewStaticCredentials(s.AccessKeyID, s.SecretAccessKey, s.SessionToken),
	})
	if err != nil {
		return err
	}
	s.cli = s3.New(sess)
	return nil
}

func (s *AmazonS3) IsZero() bool {
	return s.Region == "" &&
		s.AccessKeyID == "" &&
		s.SecretAccessKey == "" &&
		s.SessionToken == "" &&
		s.BucketName == ""
}

func (s *AmazonS3) Name() string {
	return "s3-cli"
}

func (s *AmazonS3) Upload(key string, data []byte) error {
	if _, err := s.cli.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}); err != nil {
		return err
	}

	return nil
}

func (s *AmazonS3) Read(key string) ([]byte, error) {
	resp, err := s.cli.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s *AmazonS3) Delete(key string) error {
	if _, err := s.cli.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}); err != nil {
		return err
	}

	return nil
}
