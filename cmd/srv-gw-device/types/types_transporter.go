package types

import (
	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/enums"
)

type Transporter struct {
	Mode     enums.TransportMode `env:""`
	Endpoint types.Endpoint      `env:""`
}
