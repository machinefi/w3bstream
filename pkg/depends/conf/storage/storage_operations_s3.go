package storage

import (
	"bytes"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
)

type S3 struct {
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

func (s *S3) Type() StorageType { return STORAGE_TYPE__S3 }

func (s *S3) Init() error {
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

func (s *S3) SetDefault() {
	if s.UrlExpire == 0 {
		s.UrlExpire = types.Duration(10 * time.Minute)
	}
}

func (s *S3) IsZero() bool {
	return s.Endpoint == "" ||
		s.Region == "" ||
		s.AccessKeyID == "" ||
		s.SecretAccessKey == "" ||
		s.BucketName == ""
}

func (s *S3) Name() string {
	return "s3-cli"
}

func (s *S3) Upload(key string, data []byte, chk ...HmacAlgType) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}

	t := HMAC_ALG_TYPE__MD5
	if len(chk) > 0 && chk[0] != 0 {
		t = chk[0]
	}

	sum := t.Base64Sum(data)
	switch t {
	case HMAC_ALG_TYPE__MD5:
		input.SetContentMD5(sum)
	case HMAC_ALG_TYPE__SHA1:
		input.SetChecksumAlgorithm(t.Type())
		input.SetChecksumSHA1(sum)
	case HMAC_ALG_TYPE__SHA256:
		input.SetChecksumAlgorithm(t.Type())
		input.SetChecksumSHA256(sum)
	}

	_, err := s.cli.PutObject(input)
	return err
}

func (s *S3) Read(key string, chk ...HmacAlgType) (data []byte, sum []byte, err error) {
	resp, err := s.cli.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	t := HMAC_ALG_TYPE__MD5
	if len(chk) > 0 && chk[0] != 0 {
		t = chk[0]
	}
	sum = t.Sum(data)

	return
}

func (s *S3) Delete(key string) error {
	_, err := s.cli.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3) DownloadUrl(key string) (string, error) {
	req, _ := s.cli.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
	})
	return req.Presign(s.UrlExpire.Duration())
}
