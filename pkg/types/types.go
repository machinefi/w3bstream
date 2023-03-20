package types

import "github.com/machinefi/w3bstream/pkg/depends/base/types"

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

type ETHClientConfig struct {
	Endpoints string `env:""`
}

// aliases from base/types
type (
	SFID       = types.SFID
	SFIDs      = types.SFIDs
	EthAddress = types.EthAddress
	Timestamp  = types.Timestamp
)
