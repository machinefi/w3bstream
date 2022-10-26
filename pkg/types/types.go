package types

import "github.com/iotexproject/Bumblebee/base/types"

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

type (
	SFID  = types.SFID
	SFIDs = types.SFIDs
)

var (
	EVENTTYPEDEFAULT = "DEFAULT"
)
