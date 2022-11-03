package modelgen

import "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"

type Keys struct {
	Primary       []string
	Indexes       builder.Indexes
	UniqueIndexes builder.Indexes
}

func (ks *Keys) WithSoftDelete(f string) {
	if len(ks.UniqueIndexes) > 0 {
		for name, fields := range ks.UniqueIndexes {
			ks.UniqueIndexes[name] = uniqueStrings(append(fields, f))
		}
	}
}

func (ks *Keys) Bind(tbl *builder.Table) {
	if len(ks.Primary) > 0 {
		k := &builder.Key{
			Name:     "primary",
			IsUnique: true,
			Def:      *builder.ParseIndexDef(ks.Primary...),
		}
		_ = k.Def.TableExpr(tbl)
		tbl.AddKey(k)
	}
	if len(ks.UniqueIndexes) > 0 {
		for nm, fields := range ks.UniqueIndexes {
			name, method := builder.SplitIndexNameAndMethod(nm)
			k := &builder.Key{
				Name:     name,
				Method:   method,
				IsUnique: true,
				Def:      *builder.ParseIndexDef(fields...),
			}
			_ = k.Def.TableExpr(tbl)
			tbl.AddKey(k)
		}
	}
	if len(ks.Indexes) > 0 {
		for nm, fields := range ks.Indexes {
			name, method := builder.SplitIndexNameAndMethod(nm)
			k := &builder.Key{
				Name:     name,
				Method:   method,
				IsUnique: false,
				Def:      *builder.ParseIndexDef(fields...),
			}
			_ = k.Def.TableExpr(tbl)
			tbl.AddKey(k)
		}
	}
}
