package wasmlog

import (
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type CondArgs struct {
	InstanceID types.SFID `name:"-"`
	Srcs       []string   `in:"query" name:"srcs,omitempty"`
	SrcLike    string     `in:"query" name:"src,omitempty"`
	LSrcLike   string     `in:"query" name:"lSrc,omitempty"`
}

func (r *CondArgs) Condition() builder.SqlCondition {
	var (
		m = &models.WasmLog{}
		c []builder.SqlCondition
	)
	if r.InstanceID != 0 {
		c = append(c, m.ColInstanceID().Eq(r.InstanceID))
	}
	if len(r.Srcs) > 0 {
		c = append(c, m.ColSrc().In(r.Srcs))
	}
	if r.SrcLike != "" {
		c = append(c, m.ColSrc().Like(r.SrcLike))
	}
	if r.LSrcLike != "" {
		c = append(c, m.ColSrc().LLike(r.LSrcLike))
	}
	return builder.And(c...)
}
