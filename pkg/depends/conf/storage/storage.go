package storage

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
)

var (
	ErrMissingConfigS3       = errors.New("missing config: s3")
	ErrMissingConfigFS       = errors.New("missing config: fs")
	ErrMissingConfigIPFS     = errors.New("missing config: ipfs")
	ErrUnsupprtedStorageType = errors.New("unsupported storage type")
	ErrEmptyContent          = errors.New("content empty")
	ErrContentSizeExceeded   = errors.New("content size exceeded")
	ErrDiskReservationLimit  = errors.New("disk reservation limit")
)

type StorageOperations interface {
	Type() StorageType
	Upload(key string, file []byte, chk ...HmacAlgType) error
	Read(key string, chk ...HmacAlgType) (data []byte, sum []byte, err error)
	Delete(key string) error
}

type StorageOperationsWithValidation interface {
	Validate(data []byte, sum string, chk ...HmacAlgType) bool
}

type Storage struct {
	Typ             StorageType
	FilesizeLimit   int64
	DiskReserve     int64
	PromiscuousMode bool // PromiscuousMode if support multi storage type
	TempDir         string

	*S3      `env:"S3"`
	*LocalFs `env:"Fs"`

	op StorageOperations `env:"-"`
}

func (s *Storage) Name() string { return "Storage" }

func (s *Storage) IsZero() bool {
	return s.Typ == STORAGE_TYPE_UNKNOWN || s.S3 == nil && s.LocalFs == nil
}

func (s *Storage) SetDefault() {
	if s.Typ == STORAGE_TYPE_UNKNOWN {
		s.Typ = STORAGE_TYPE__FILESYSTEM
	}
	if s.FilesizeLimit == 0 {
		s.FilesizeLimit = 1024 * 1024
	}
	if s.DiskReserve == 0 {
		s.DiskReserve = 20 * 1024 * 1024
	}
}

func (s *Storage) Init() error {
	if s.TempDir == "" {
		tmp := os.Getenv("TMPDIR")
		if tmp == "" {
			tmp = "/tmp"
		}
		service := os.Getenv(consts.EnvProjectName)
		if service == "" {
			service = "service"
		}
		s.TempDir = filepath.Join(tmp, service)
	}
	// overwrite default 'TMPDIR'
	if err := os.Setenv("TMPDIR", s.TempDir); err != nil {
		return err
	}

	switch s.Typ {
	case STORAGE_TYPE_UNKNOWN, STORAGE_TYPE__FILESYSTEM:
		if s.LocalFs == nil {
			return ErrMissingConfigFS
		}
		s.op = s.LocalFs
		s.Typ = STORAGE_TYPE__FILESYSTEM
	case STORAGE_TYPE__S3:
		if s.S3 == nil || s.S3.IsZero() {
			return ErrMissingConfigS3
		}
		s.op = s.S3
	case STORAGE_TYPE__IPFS:
		return ErrMissingConfigIPFS
	default:
		return ErrUnsupprtedStorageType
	}

	if canSetDefault, ok := s.op.(types.DefaultSetter); ok {
		canSetDefault.SetDefault()
	}
	if canBeInit, ok := s.op.(types.Initializer); ok {
		canBeInit.Init()
		return nil
	}
	if canBeInit, ok := s.op.(types.ValidatedInitializer); ok {
		return canBeInit.Init()
	}
	return nil
}

func (s *Storage) WithOperation(op StorageOperations) {
	s.op = op
}

func (s *Storage) Upload(key string, content []byte, chk ...HmacAlgType) error {
	size := int64(len(content))
	if size == 0 {
		return ErrEmptyContent
	}
	if s.FilesizeLimit != 0 && size > s.FilesizeLimit {
		return ErrContentSizeExceeded
	}

	free := int64(0)
	stat, err := disk.Usage(s.TempDir)
	if err == nil {
		free = int64(stat.Free) - size
	}

	if s.DiskReserve != 0 && s.Type() == STORAGE_TYPE__FILESYSTEM && free < s.DiskReserve {
		return ErrDiskReservationLimit
	}

	if err = s.op.Upload(key, content, chk...); err != nil {
		return err
	}
	if s.Typ == STORAGE_TYPE__S3 {
		_ = os.RemoveAll(s.TempDir)
	}
	return nil
}

func (s *Storage) Read(key string, chk ...HmacAlgType) ([]byte, []byte, error) {
	return s.op.Read(key, chk...)
}

func (s *Storage) Delete(key string) error {
	return s.op.Delete(key)
}

func (s *Storage) Validate(data []byte, sum string, chk ...HmacAlgType) bool {
	if len(data) == 0 || len(sum) == 0 {
		return true
	}

	t := HMAC_ALG_TYPE__MD5
	if len(chk) > 0 && chk[0] != 0 {
		t = chk[0]
	}

	return sum == t.HexSum(data)
}

func (s *Storage) Type() StorageType { return s.op.Type() }
