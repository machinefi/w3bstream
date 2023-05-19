package amazonS3

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
)

type AmazonS3 struct {
	Endpoint         string         `env:""`
	Region           string         `env:""`
	AccessKeyID      string         `env:""`
	SecretAccessKey  types.Password `env:""`
	SessionToken     string         `env:""`
	BucketName       string         `env:""`
	S3ForcePathStyle bool           `env:""`

	cli *s3.S3
}

func (s *AmazonS3) Init() error {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(s.Endpoint),
		Region:           aws.String(s.Region),
		Credentials:      credentials.NewStaticCredentials(s.AccessKeyID, s.SecretAccessKey.String(), s.SessionToken),
		S3ForcePathStyle: aws.Bool(s.S3ForcePathStyle),
	})
	if err != nil {
		return err
	}
	s.cli = s3.New(sess)
	return nil
}

func (s *AmazonS3) IsZero() bool {
	return s.Endpoint == "" ||
		s.Region == "" ||
		s.AccessKeyID == "" ||
		s.SecretAccessKey == "" ||
		s.BucketName == ""
}

func (s *AmazonS3) Name() string {
	return "s3-cli"
}

func (s *AmazonS3) Upload(key string, data []byte) error {
	_, err := s.cli.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	return err
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
	_, err := s.cli.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	return err
}
