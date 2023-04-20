package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
)

type S3 struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	BucketName      string

	cli *s3.S3
}

func (s *S3) SetDefault() {
	if s.Region == "" {
		s.Region = "us-west-1"
	}
}

func (s *S3) Init() error {
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

func (s *S3) Name() string {
	return "s3-cli"
}

//func (s *S3) LivenessCheck() map[string]string {
//	return nil
//}

func (s *S3) Upload(key string, file []byte) error {
	if _, err := s.cli.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	}); err != nil {
		return err
	}

	return nil
}

func (s *S3) Read(key string) ([]byte, error) {
	resp, err := s.cli.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *S3) Delete(key string) error {
	if _, err := s.cli.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	}); err != nil {
		return err
	}

	return nil
}
