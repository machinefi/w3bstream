package storage_test

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/storage"
	mock_conf_storage "github.com/machinefi/w3bstream/pkg/test/mock_depends_conf_storage"
)

func TestStorage(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	t.Run("IsZero", func(t *testing.T) {
		s := &storage.Storage{Typ: storage.STORAGE_TYPE_UNKNOWN}
		NewWithT(t).Expect(s.IsZero()).To(BeTrue())

		s = &storage.Storage{
			Typ: storage.STORAGE_TYPE__S3,
			S3:  &storage.S3{},
		}
		NewWithT(t).Expect(s.IsZero()).To(BeFalse())
	})

	t.Run("SetDefault", func(t *testing.T) {
		s := &storage.Storage{Typ: storage.STORAGE_TYPE_UNKNOWN}
		s.SetDefault()
		NewWithT(t).Expect(s.Typ).To(Equal(storage.STORAGE_TYPE__FILESYSTEM))

		s = &storage.Storage{}
		s.SetDefault()
		NewWithT(t).Expect(s.FilesizeLimit).To(Equal(int64(1024 * 1024)))
		NewWithT(t).Expect(s.DiskReserve).To(Equal(int64(20 * 1024 * 1024)))

		s = &storage.Storage{
			FilesizeLimit: 100,
			DiskReserve:   100,
		}
		s.SetDefault()
		NewWithT(t).Expect(s.FilesizeLimit).To(Equal(int64(100)))
		NewWithT(t).Expect(s.DiskReserve).To(Equal(int64(100)))
	})

	t.Run("Init", func(t *testing.T) {
		t.Run("#InitTempDir", func(t *testing.T) {
			s := &storage.Storage{LocalFs: &storage.LocalFs{}}
			cases := []*struct {
				preFn  func()
				expect string
			}{
				{
					preFn: func() {
						_ = os.Unsetenv("TMPDIR")
						_ = os.Unsetenv(consts.EnvProjectName)
					},
					expect: "/tmp/service",
				},
				{
					preFn: func() {
						_ = os.Setenv("TMPDIR", "/test_tmp")
						_ = os.Setenv(consts.EnvProjectName, "test_storage")
					},
					expect: "/test_tmp/test_storage",
				},
			}

			for _, cc := range cases {
				cc.preFn()
				err := s.Init()
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(s.TempDir).To(Equal(os.Getenv("TMPDIR")))
			}
		})

		t.Run("#InitTypeAndOp", func(t *testing.T) {
			cases := []*struct {
				conf   *storage.Storage
				expect error
			}{{
				conf:   &storage.Storage{},
				expect: storage.ErrMissingConfigFS,
			}, {
				conf:   &storage.Storage{LocalFs: &storage.LocalFs{}},
				expect: nil,
			}, {
				conf:   &storage.Storage{Typ: storage.STORAGE_TYPE__S3},
				expect: storage.ErrMissingConfigS3,
			}, {
				conf: &storage.Storage{
					Typ: storage.STORAGE_TYPE__S3,
					S3: &storage.S3{
						Endpoint:        "http://demo.s3.org",
						Region:          "us",
						AccessKeyID:     "1",
						SecretAccessKey: "1",
						BucketName:      "test_bucket",
					},
				},
				expect: nil,
			}, {
				conf:   &storage.Storage{Typ: storage.STORAGE_TYPE__IPFS},
				expect: storage.ErrMissingConfigIPFS,
			}, {
				conf:   &storage.Storage{Typ: storage.StorageType(100)},
				expect: storage.ErrUnsupprtedStorageType,
			}}

			for idx, cc := range cases {
				t.Run("#"+strconv.Itoa(idx), func(t *testing.T) {
					err := cc.conf.Init()
					if cc.expect == nil {
						NewWithT(t).Expect(err).To(BeNil())
					} else {
						NewWithT(t).Expect(err).To(Equal(cc.expect))
					}
				})
			}
		})
	})

	t.Run("#Upload", func(t *testing.T) {
		cc := &storage.Storage{TempDir: "/tmp"}

		t.Run("#Success", func(t *testing.T) {
			op := mock_conf_storage.NewMockStorageOperations(c)
			op.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
			cc.WithOperation(op)

			err := cc.Upload("any", []byte("any"))
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#EmptyContent", func(t *testing.T) {
				err := cc.Upload("any", []byte(""))
				NewWithT(t).Expect(err).To(Equal(storage.ErrEmptyContent))
			})
			t.Run("#FileSizeLimit", func(t *testing.T) {
				cc.FilesizeLimit = 4
				err := cc.Upload("any", []byte("12345"))
				NewWithT(t).Expect(err).To(Equal(storage.ErrContentSizeExceeded))
			})
			t.Run("#DiskReserve", func(t *testing.T) {
				op := mock_conf_storage.NewMockStorageOperations(c)

				op.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(1)
				op.EXPECT().Type().Return(storage.STORAGE_TYPE__FILESYSTEM).MaxTimes(1)
				cc.WithOperation(op)

				stat, err := disk.Usage(cc.TempDir)
				NewWithT(t).Expect(err).To(BeNil())

				cc.DiskReserve = int64(stat.Free + 1024*1024*1024)
				cc.FilesizeLimit = 0

				err = cc.Upload("any", []byte("any"))
				NewWithT(t).Expect(err).To(Equal(storage.ErrDiskReservationLimit))
			})
			t.Run("#OpUploadFailed", func(t *testing.T) {
				op := mock_conf_storage.NewMockStorageOperations(c)

				op.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(errors.New("mock error")).MaxTimes(1)
				cc.WithOperation(op)

				cc.DiskReserve = 0
				cc.FilesizeLimit = 0

				err := cc.Upload("any", []byte("any"))
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})

	t.Run("#Read", func(t *testing.T) {
		cc := &storage.Storage{}

		t.Run("#Success", func(t *testing.T) {
			op := mock_conf_storage.NewMockStorageOperations(c)
			op.EXPECT().Read(gomock.Any()).Return(nil, nil, nil).MaxTimes(1)
			cc.WithOperation(op)

			_, _, err := cc.Read("any")
			NewWithT(t).Expect(err).To(BeNil())
		})
		t.Run("#Failed", func(t *testing.T) {
			op := mock_conf_storage.NewMockStorageOperations(c)
			op.EXPECT().Read(gomock.Any()).Return(nil, nil, errors.New("mock error")).MaxTimes(1)
			cc.WithOperation(op)

			_, _, err := cc.Read("any")
			NewWithT(t).Expect(err).NotTo(BeNil())
		})
	})

	t.Run("#Type", func(t *testing.T) {
		cc := &storage.Storage{}
		expect := storage.STORAGE_TYPE__S3

		op := mock_conf_storage.NewMockStorageOperations(c)
		op.EXPECT().Type().Return(expect).MaxTimes(1)
		cc.WithOperation(op)

		NewWithT(t).Expect(cc.Type()).To(Equal(expect))
	})

	t.Run("#Validate", func(t *testing.T) {
		cc := &storage.Storage{}

		content := []byte("1234567")
		md5sum := storage.HMAC_ALG_TYPE__MD5.HexSum(content)
		sha1sum := storage.HMAC_ALG_TYPE__SHA1.HexSum(content)
		sha256sum := storage.HMAC_ALG_TYPE__SHA256.HexSum(content)

		NewWithT(t).Expect(cc.Validate(nil, "sum")).To(BeTrue())
		NewWithT(t).Expect(cc.Validate([]byte("xx"), "")).To(BeTrue())
		NewWithT(t).Expect(cc.Validate(content, md5sum)).To(BeTrue())
		NewWithT(t).Expect(cc.Validate(content, sha1sum, storage.HMAC_ALG_TYPE__SHA1)).To(BeTrue())
		NewWithT(t).Expect(cc.Validate(content, sha256sum, storage.HMAC_ALG_TYPE__SHA256)).To(BeTrue())
	})
}

func TestS3(t *testing.T) {
	t.Run("IsZero", func(t *testing.T) {
		var (
			valued = &storage.S3{
				Endpoint:        "s3://sincos-test",
				Region:          "us-east-2",
				AccessKeyID:     "xx",
				SecretAccessKey: "xx",
				BucketName:      "sincos-test",
			}
			empty = &storage.S3{}
		)
		NewWithT(t).Expect(valued.IsZero()).To(BeFalse())
		NewWithT(t).Expect(empty.IsZero()).To(BeTrue())
	})

	t.Run("SetDefault", func(t *testing.T) {
		var (
			dftExpiration = types.Duration(10 * time.Minute)
			conf          = &storage.S3{UrlExpire: dftExpiration / 2}
		)
		conf.UrlExpire = 0
		conf.SetDefault()
		NewWithT(t).Expect(conf.UrlExpire).To(Equal(dftExpiration))
		conf.UrlExpire = dftExpiration / 2
		conf.SetDefault()
		NewWithT(t).Expect(conf.UrlExpire).To(Equal(dftExpiration / 2))
	})

	t.Run("Init", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			conf := &storage.S3{
				Endpoint:   "",
				Region:     "us-east-2",
				BucketName: "test",
			}
			err := conf.Init()
			NewWithT(t).Expect(err).To(BeNil())
		})

		t.Run("#Failed", func(t *testing.T) {
			t.Run("#NewSessionFailed", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}
				patcher := gomonkey.ApplyFunc(
					session.NewSession,
					func(...*aws.Config) (*session.Session, error) {
						return nil, errors.New("")
					},
				)
				defer patcher.Reset()

				conf := &storage.S3{}
				err := conf.Init()
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})

	var (
		// ep = &storage.S3{}
		ep = &storage.S3{
			Endpoint:         "s3.us-west-1.amazonaws.com",
			Region:           "us-west-1",
			AccessKeyID:      "AKIAVJEGDD6ZONM3RO6V",
			SecretAccessKey:  "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",
			BucketName:       "machinefi-w3bstream-devnet-staging",
			S3ForcePathStyle: true,
		}
		key       = "unit_test_key"
		data      = []byte("unit_test_data")
		sumMd5    = storage.HMAC_ALG_TYPE__MD5.Sum(data)
		sumSha1   = storage.HMAC_ALG_TYPE__SHA1.Sum(data)
		sumSha256 = storage.HMAC_ALG_TYPE__SHA256.Sum(data)
	)

	t.Run("Upload", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			patch := gomonkey.ApplyMethod(
				reflect.TypeOf(&s3.S3{}),
				"PutObject",
				func(receiver *s3.S3, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return nil, nil
				},
			)
			defer patch.Reset()

			t.Run("#WithoutSumCheck", func(t *testing.T) {
				err := ep.Upload(key, data)
				NewWithT(t).Expect(err).To(BeNil())
			})
			t.Run("#WithMd5SumCheck", func(t *testing.T) {
				err := ep.Upload(key, data, storage.HMAC_ALG_TYPE__MD5)
				NewWithT(t).Expect(err).To(BeNil())
			})
			t.Run("#WithSHA1SumCheck", func(t *testing.T) {
				err := ep.Upload(key, data, storage.HMAC_ALG_TYPE__SHA1)
				NewWithT(t).Expect(err).To(BeNil())
			})
			t.Run("#WithSHA256SumCheck", func(t *testing.T) {
				err := ep.Upload(key, data, storage.HMAC_ALG_TYPE__SHA256)
				NewWithT(t).Expect(err).To(BeNil())
			})
		})
		t.Run("#Failed", func(t *testing.T) {
			patch := gomonkey.ApplyMethod(
				reflect.TypeOf(&s3.S3{}),
				"PutObject",
				func(receiver *s3.S3, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
					return nil, errors.New("")
				},
			)
			defer patch.Reset()
			err := ep.Upload(key, data)
			NewWithT(t).Expect(err).NotTo(BeNil())
		})
	})
	t.Run("Read", func(t *testing.T) {
		t.Run("#Success", func(t *testing.T) {
			patch := gomonkey.ApplyMethod(
				reflect.TypeOf(&s3.S3{}),
				"GetObject",
				func(recv *s3.S3, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
					return &s3.GetObjectOutput{
						Body: io.NopCloser(bytes.NewBuffer(data)),
					}, nil
				},
			)
			defer patch.Reset()

			cases := []*struct {
				chk []storage.HmacAlgType
				sum []byte
			}{
				{chk: nil, sum: sumMd5},
				{chk: []storage.HmacAlgType{storage.HMAC_ALG_TYPE__MD5}, sum: sumMd5},
				{chk: []storage.HmacAlgType{storage.HMAC_ALG_TYPE__SHA1}, sum: sumSha1},
				{chk: []storage.HmacAlgType{storage.HMAC_ALG_TYPE__SHA256}, sum: sumSha256},
			}

			for _, c := range cases {
				content, sum, err := ep.Read(key, c.chk...)
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(bytes.Equal(content, data)).To(BeTrue())
				NewWithT(t).Expect(bytes.Equal(sum, c.sum)).To(BeTrue())
			}
		})
		t.Run("#Failed", func(t *testing.T) {
		})
	})
	t.Run("Delete", func(t *testing.T) {})
}
