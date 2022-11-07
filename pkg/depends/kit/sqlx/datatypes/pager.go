package datatypes

import "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"

type Pager struct {
	Size   int64 `name:"size,omitempty"   in:"query" default:"10" validate:"@int64[-1,]"`
	Offset int64 `name:"offset,omitempty" in:"query" default:"0"  validate:"@int64[0,]"`
}

func (p *Pager) Addition() builder.Addition {
	if p.Size != -1 {
		return builder.Limit(p.Size).Offset(p.Offset)
	}
	return nil
}
