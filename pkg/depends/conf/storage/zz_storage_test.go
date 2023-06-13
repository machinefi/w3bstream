package storage_test

import (
	"os"
	"runtime"
	"strconv"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
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

		t.Run("#Upload", func(t *testing.T) {
			cc := &storage.Storage{TempDir: "/tmp"}
			op := mock_conf_storage.NewMockStorageOperations(c)

			t.Run("#Success", func(t *testing.T) {
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
				t.Run("#FileSizeLimit", func(t *testing.T) {
					if runtime.GOOS == `darwin` {
						return
					}
					cc.FilesizeLimit = 0
					cc.DiskReserve = 100

					patches := gomonkey.ApplyFunc(
						disk.Usage,
						func(_ string) (*disk.UsageStat, error) {
							return &disk.UsageStat{Free: 1}, nil
						},
					)
					defer patches.Reset()

					err := cc.Upload("any", []byte("any"))
					NewWithT(t).Expect(err).To(Equal(storage.ErrDiskReservationLimit))
				})
			})
		})
	})
}

/*

	cc.WithOperation(op)

	op.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(errors.New("")).MaxTimes(1)
	NewWithT()

	op.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(errors.New("")).MaxTimes(1)

*/
