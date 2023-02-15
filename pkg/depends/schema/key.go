package schema

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
)

type Keys struct {
	lst []*Key
	*mapx.Map[string, *Key]
}

func (ks *Keys) Add(keys ...*Key) {
	if ks.Map == nil {
		ks.Map = mapx.New[string, *Key]()
	}
	for _, k := range keys {
		if k != nil {
			ks.lst = append(ks.lst, k)
			if !ks.StoreNX(strings.ToLower(k.Name), k) {
				panic(errors.Errorf("duplicated key: %s", k.Name))
			}
		}
	}
}

func (ks *Keys) Range(f func(k *Key, idx int)) {
	for idx, k := range ks.lst {
		f(k, idx)
	}
}

func (ks *Keys) Reset() {
	ks.lst = ks.lst[0:0]
	if ks.Map != nil {
		ks.Clear()
	}
}

type Key struct {
	Name     string `json:"name,omitempty"`     // Name index name
	Method   string `json:"method,omitempty"`   // Method index method `btree` or `hash`
	IsUnique bool   `json:"isUnique,omitempty"` // IsUnique is unique index
	IndexDef
	WithTableDefinition `json:"-"`
}

func (k *Key) IsPrimary() bool {
	return k.IsUnique && (k.Name == "primary" || strings.HasPrefix(k.Name, "pkey"))
}

type IndexDef struct {
	ColumnNames []string `json:"columnNames"`    // ColumnNames index related column names
	Expr        string   `json:"expr,omitempty"` // Expr index extension for sql dialect
}
