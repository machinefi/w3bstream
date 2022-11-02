package types

import "github.com/machinefi/Bumblebee/base/types"

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
	PrivateKey    string `env:"" yaml:"privateKey"`
	ChainEndpoint string `env:"" yaml:"chainEndpoint"`
}

type (
	SFID  = types.SFID
	SFIDs = types.SFIDs
)
