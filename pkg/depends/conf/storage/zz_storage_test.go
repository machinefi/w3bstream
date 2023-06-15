package storage_test

import (
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/conf/storage"
	"github.com/machinefi/w3bstream/pkg/test/mock_depends_conf_storage"
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
	conf := &storage.S3{
		Endpoint:        "s3://sincos-test",
		Region:          "us-east-2",
		AccessKeyID:     "xx",
		SecretAccessKey: "xx",
		BucketName:      "sincos-test",
	}
	t.Run("IsZero", func(t *testing.T) {
		var (
			valued = &(*conf)
			empty  = &storage.S3{}
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
		t.Run("#Failed", func(t *testing.T) {
			t.Run("#NewSessionFailed", func(t *testing.T) {
				if runtime.GOOS == `darwin` {
					return
				}
				patcher := gomonkey.ApplyFunc(
					session.NewSession,
					func(...*aws.Config) (*session.Session, error) {
						return &session.Session{}, errors.New("")
					},
				)
				defer patcher.Reset()

				err := conf.Init()
				NewWithT(t).Expect(err).NotTo(BeNil())
			})
		})
	})
	t.Run("Upload", func(t *testing.T) {})
	t.Run("Read", func(t *testing.T) {})
	t.Run("Delete", func(t *testing.T) {})
}
