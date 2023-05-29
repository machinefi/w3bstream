package amazonS3

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"

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
	return s.UploadWithMD5(key, "", data)
}

func (s *AmazonS3) UploadWithMD5(key, md5 string, data []byte) error {
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	if md5 != "" {
		md5Bytes, err := hex.DecodeString(md5)
		if err != nil {
			return err
		}
		putObjectInput.SetContentMD5(base64.StdEncoding.EncodeToString(md5Bytes))
	}

	_, err := s.cli.PutObject(putObjectInput)
	return err
}

func (s *AmazonS3) Read(key string) ([]byte, error) {
	return s.ReadWithMD5(key, "")
}

func (s *AmazonS3) ReadWithMD5(key, md5 string) ([]byte, error) {
	resp, err := s.cli.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if md5 != "" {
		sum := strings.Trim(*resp.ETag, "\"")
		if sum != md5 {
			return nil, errors.New("md5 not match")
		}
	}

	return io.ReadAll(resp.Body)
}

func (s *AmazonS3) Delete(key string) error {
	_, err := s.cli.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	return err
}
