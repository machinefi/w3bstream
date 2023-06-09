package storage_test

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/conf/storage"
)

func TestStorage_Init(t *testing.T) {
	cases := []struct {
		conf                  *storage.Storage
		expectedInitErr       error
		setOperations         func(c *storage.Storage)
		expectedErrAfterSetOp error
	}{
		{
			conf: &storage.Storage{
				Type: storage.STORAGE_TYPE__S3,
				S3:   &storage.S3{},
			},
			expectedInitErr: storage.ErrMissingConfigS3,
			setOperations: func(conf *storage.Storage) {
				conf.S3 = &storage.S3{
					Endpoint:        "http://1.1.1.1",
					Region:          "us",
					AccessKeyID:     "111",
					SecretAccessKey: "222",
					SessionToken:    "",
					BucketName:      "endpoint_test",
				}
			},
		}, {
			conf: &storage.Storage{
				Type:    storage.STORAGE_TYPE__FILESYSTEM,
				LocalFs: &storage.LocalFs{},
			},
		}, {
			conf: &storage.Storage{
				Type: storage.STORAGE_TYPE__IPFS,
			},
			expectedInitErr: storage.ErrMissingConfigIPFS,
		}, {
			conf: &storage.Storage{
				Type: storage.StorageType(100),
			},
			expectedInitErr: storage.ErrUnsupprtedStorageType,
		},
	}

	for _, c := range cases {
		t.Run(c.conf.Type.String(), func(t *testing.T) {
			c.conf.SetDefault()
			NewWithT(t).Expect(c.conf.DiskReserve).To(Equal(int64(20 * 1024 * 1024)))

			err := c.conf.Init()
			if err != nil {
				NewWithT(t).Expect(err.Error()).To(Equal(c.expectedInitErr.Error()))
			}

			if c.setOperations != nil {
				c.setOperations(c.conf)
				if err = c.conf.Init(); err != nil {
					NewWithT(t).Expect(err).To(Equal(c.expectedErrAfterSetOp))
				}
			}

			NewWithT(t).Expect(os.Getenv("TMPDIR")).To(Equal(c.conf.TempDir))
		})
	}
}
