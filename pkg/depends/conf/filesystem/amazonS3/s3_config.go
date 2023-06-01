package amazonS3

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/filesystem"
	"github.com/machinefi/w3bstream/pkg/depends/util"
)

type AmazonS3 struct {
	Endpoint         string         `env:""`
	Region           string         `env:""`
	AccessKeyID      string         `env:""`
	SecretAccessKey  types.Password `env:""`
	SessionToken     string         `env:""`
	BucketName       string         `env:""`
	UrlExpire        types.Duration `env:""`
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

func (s *AmazonS3) SetDefault() {
	if s.UrlExpire == 0 {
		s.UrlExpire = types.Duration(10 * time.Minute)
	}
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
	return s.UploadWithChecksum(key, "", "", data)
}

func (s *AmazonS3) UploadWithChecksum(key, sum, algorithm string, data []byte) error {
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	if sum != "" && algorithm != "" {
		sumBytes, err := hex.DecodeString(sum)
		if err != nil {
			return err
		}
		switch algorithm {
		case util.Md5Algorithm:
			putObjectInput.SetContentMD5(base64.StdEncoding.EncodeToString(sumBytes))
		case util.Sha256Algorithm:
			putObjectInput.SetChecksumAlgorithm(util.Sha256Algorithm)
			putObjectInput.SetChecksumSHA256(base64.StdEncoding.EncodeToString(sumBytes))
		}
	}

	_, err := s.cli.PutObject(putObjectInput)
	return err
}

func (s *AmazonS3) Read(key string) ([]byte, error) {
	return s.ReadWithChecksum(key, "", "")
}

func (s *AmazonS3) ReadWithChecksum(key, sum, algorithm string) ([]byte, error) {
	resp, err := s.cli.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if sum != "" && algorithm != "" {
		var targetSum string
		switch algorithm {
		case util.Md5Algorithm:
			targetSum = strings.Trim(*resp.ETag, "\"")
		case util.Sha256Algorithm:
			targetSum = *resp.ChecksumSHA256
		}
		if sum != targetSum {
			return nil, filesystem.ErrChecksumNotMatch
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

func (s *AmazonS3) DownloadUrl(key string) (string, error) {
	req, _ := s.cli.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	return req.Presign(s.UrlExpire.Duration())
}
