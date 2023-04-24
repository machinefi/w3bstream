package config

import (
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

type CondArgs struct {
	ConfigIDs []types.SFID
	RelIDs    []types.SFID
	Types     []enums.ConfigType
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m = &models.Config{}
		c []builder.SqlCondition
	)

	if len(r.ConfigIDs) > 0 {
		c = append(c, m.ColConfigID().In(r.ConfigIDs))
	}
	if len(r.RelIDs) > 0 {
		c = append(c, m.ColRelID().In(r.RelIDs))
	}
	if len(r.Types) > 0 {
		c = append(c, m.ColType().In(r.Types))
	}
	return builder.And(c...)
}

type Detail struct {
	RelID types.SFID
	wasm.Configuration
}

func (d *Detail) String() string {
	return fmt.Sprintf("[rel: %v][type: %v]", d.RelID, d.ConfigType())
}

func (d *Detail) Log(err error) string {
	s := d.String()
	if err == nil {
		return s
	}
	return fmt.Sprintf("%s: %v", s, err)
}
