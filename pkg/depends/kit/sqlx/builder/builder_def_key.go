package builder

import "strings"

type Key struct {
	Table    *Table
	Name     string
	IsUnique bool
	Method   string
	Def      IndexDef
}

func (k Key) On(t *Table) *Key { k.Table = t; return &k }

func (k Key) Using(method string) *Key { k.Method = method; return &k }

func (k *Key) T() *Table { return k.Table }

func (k Key) IsPrimary() bool {
	return k.IsUnique && (k.Name == "primary" || strings.HasSuffix(k.Name, "pkey"))
}

type Keys struct {
	lst []*Key
}

func (ks *Keys) Len() int {
	if ks == nil {
		return 0
	}
	return len(ks.lst)
}

func (ks *Keys) Clone() *Keys {
	cloned := &Keys{}
	ks.Range(func(k *Key, idx int) {
		cloned.Add(k)
	})
	return cloned
}

func (ks *Keys) Range(f func(k *Key, idx int)) {
	for i := range ks.lst {
		f(ks.lst[i], i)
	}
}

func (ks *Keys) Add(keys ...*Key) {
	for i := range keys {
		if k := keys[i]; k != nil {
			ks.lst = append(ks.lst, k)
		}
	}
}

func (ks *Keys) Key(name string) *Key {
	name = strings.ToLower(name)
	for i := range ks.lst {
		if name == ks.lst[i].Name {
			return ks.lst[i]
		}
	}
	return nil
}

// IndexDef @def index xxx/BTREE FieldA FieldB ...
type IndexDef struct {
	FieldNames []string
	ColNames   []string
	Expr       string
}

func ParseIndexDef(names ...string) *IndexDef {
	f := &IndexDef{}

	if len(names) == 1 {
		s := names[0]
		if strings.Contains(s, "#") || strings.Contains(s, "(") {
			f.Expr = s
		} else {
			f.FieldNames = strings.Split(s, " ")
		}
	} else {
		f.FieldNames = names
	}
	return f
}

func (i IndexDef) ToDefs() []string {
	if i.Expr != "" {
		return []string{i.Expr}
	}
	return i.FieldNames
}

func (i IndexDef) TableExpr(t *Table) *Ex {
	if i.Expr != "" {
		return t.Expr(i.Expr)
	}
	if len(i.ColNames) != 0 {
		ex := Expr("")
		ex.WriteGroup(func(ex *Ex) {
			ex.WriteExpr(t.MustCols(i.ColNames...))
		})
		return ex
	}
	ex := Expr("")
	ex.WriteGroup(func(ex *Ex) {
		ex.WriteExpr(t.MustColsByFieldNames(i.FieldNames...))
	})
	return ex
}

func PrimaryKey(cols *Columns) *Key { return UniqueIndex("PRIMARY", cols) }

func UniqueIndex(name string, cols *Columns, exprs ...string) *Key {
	k := Index(name, cols, exprs...)
	k.IsUnique = true
	return k
}

func Index(name string, cols *Columns, exprs ...string) *Key {
	k := &Key{Name: strings.ToLower(name)}
	if cols != nil {
		k.Def.FieldNames = cols.FieldNames()
		k.Def.ColNames = cols.ColNames()
	}
	if len(exprs) > 0 {
		k.Def.Expr = strings.Join(exprs, " ")
	}
	return k
}
