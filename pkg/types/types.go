package types

import (
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/strfmt"
	"github.com/machinefi/w3bstream/pkg/enums"
)

type UploadConfig struct {
	Root          string `env:""`
	FileSizeLimit int64  `env:""`
}

func (c *UploadConfig) SetDefault() {
	if c.Root == "" {
		c.Root = "./asserts"
	}
	if c.FileSizeLimit == 0 {
		c.FileSizeLimit = 100 * 1024 * 1024
	}
}

type FileSystem struct {
	Type enums.FileSystemMode `env:""`
}

func (f *FileSystem) SetDefault() {
	if f.Type > enums.FILE_SYSTEM_MODE__S3 || f.Type <= 0 {
		f.Type = enums.FILE_SYSTEM_MODE__LOCAL
	}
}

type ETHClientConfig struct {
	Endpoints string `env:""`
}

// aliases from base/types
type (
	SFID                     = types.SFID
	SFIDs                    = types.SFIDs
	EthAddress               = types.EthAddress
	Timestamp                = types.Timestamp
	Initializer              = types.Initializer
	ValidatedInitializer     = types.ValidatedInitializer
	InitializerWith          = types.InitializerWith
	ValidatedInitializerWith = types.ValidatedInitializerWith
)

type WhiteList []string

func (v *WhiteList) Init() {
	lst := WhiteList{}
	for _, addr := range *v {
		if err := strfmt.EthAddressValidator.Validate(addr); err == nil {
			lst = append(lst, strings.ToLower(addr))
		}
	}
	*v = lst
}

func (v *WhiteList) Validate(address string) bool {
	if v == nil || len(*v) == 0 {
		return true
	}
	for _, addr := range *v {
		if addr == strings.ToLower(address) {
			return true
		}
	}
	return false
}

type StrategyResult struct {
	AccountID   types.SFID `json:"accountID"   db:"f_acc_id"`
	ProjectID   types.SFID `json:"projectID"   db:"f_prj_id"`
	ProjectName string     `json:"projectName" db:"f_prj_name"`
	AppletID    types.SFID `json:"appletID"    db:"f_app_id"`
	AppletName  string     `json:"appletName"  db:"f_app_name"`
	InstanceID  types.SFID `json:"instanceID"  db:"f_ins_id"`
	Handler     string     `json:"handler"     db:"f_hdl"`
	EventType   string     `json:"eventType"   db:"f_evt"`
}
