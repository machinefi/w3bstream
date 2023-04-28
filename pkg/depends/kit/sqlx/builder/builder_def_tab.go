package builder

import (
	"bytes"
	"container/list"
	"context"
	"fmt"
	"sort"
	"strings"
	"text/scanner"
)

type TableDefinition interface{ T() *Table }

func T(tblName string, defs ...TableDefinition) *Table {
	t := &Table{Name: tblName}

	for _, def := range defs {
		if d, ok := def.(*Column); ok {
			t.AddCol(d)
		}
		if k, ok := def.(*Key); ok {
			t.AddKey(k)
		}
	}
	return t
}

type Table struct {
	Name      string
	Desc      []string
	Schema    string
	ModelName string
	Model     Model

	Columns
	Keys
}

func (t *Table) TableName() string { return t.Name }

func (t *Table) IsNil() bool { return t == nil || t.Name == "" }

func (t *Table) Ex(ctx context.Context) *Ex {
	if t.Schema != "" {
		return Expr(t.Schema + "." + t.Name).Ex(ctx)
	}
	return Expr(t.Name).Ex(ctx)
}

func (t Table) WithSchema(schema string) *Table {
	t.Schema = schema

	cols := Columns{}
	t.Columns.Range(func(c *Column, idx int) {
		cols.Add(c.On(&t))
	})
	t.Columns = cols

	keys := Keys{}
	t.Keys.Range(func(k *Key, idx int) {
		keys.Add(k.On(&t))
	})
	t.Keys = keys
	return &t
}

func (t *Table) AddCol(c *Column) {
	if c != nil {
		t.Columns.Add(c.On(t))
	}
}

func (t *Table) AddKey(k *Key) {
	if k != nil {
		t.Keys.Add(k.On(t))
	}
}

func (t *Table) Expr(query string, args ...interface{}) *Ex {
	if query == "" {
		return nil
	}

	argc := len(args)
	e := Expr("")
	e.Grow(argc)

	s := scanner.Scanner{}
	s.Init(bytes.NewBuffer([]byte(query)))

	qc := 0

	for tok := s.Next(); tok != scanner.EOF; tok = s.Next() {
		switch tok {
		default:
			e.WriteQueryRune(tok)
		case '?':
			e.WriteQueryRune(tok)
			if qc < argc {
				e.AppendArgs(args[qc])
				qc++
			}
		case '#':
			b := bytes.NewBuffer(nil) //field name buffer

			e.WriteHolder(0)
			for {
				tok = s.Next()
				if tok == scanner.EOF {
					break
				}
				if tok >= 'A' && tok <= 'Z' || tok >= 'a' && tok <= 'z' ||
					tok >= '0' && tok <= '9' || tok == '_' {
					b.WriteRune(tok)
					continue
				}
				e.WriteQueryRune(tok)
				break
			}
			if b.Len() == 0 {
				e.AppendArgs(t)
				continue
			}
			name := b.String()
			col := t.ColByFieldName(name)
			if col == nil {
				panic(fmt.Errorf("missing field %s of %s", name, t.Name))
			}
			e.AppendArgs(col)
		}
	}
	return e
}

func (t *Table) Diff(prevT *Table, d Dialect) (exprList []SqlExpr) {
	// diff columns
	t.Columns.Range(func(currC *Column, idx int) {
		if prevC := prevT.Col(currC.Name); prevC != nil {
			if currC != nil {
				if currC.DeprecatedActs != nil {
					renameTo := currC.DeprecatedActs.RenameTo
					if renameTo != "" {
						prevCol := prevT.Col(renameTo)
						if prevCol != nil {
							exprList = append(exprList, d.DropColumn(prevCol))
						}
						targetCol := t.Col(renameTo)
						if targetCol == nil {
							panic(fmt.Errorf("col `%s` is not declared", renameTo))
						}
						exprList = append(exprList, d.RenameColumn(currC, targetCol))
						prevT.AddCol(targetCol)
						return
					}
					exprList = append(exprList, d.DropColumn(currC))
					return
				}

				prevCT := d.DataType(prevC.ColumnType).Ex(context.Background()).Query()
				currCT := d.DataType(currC.ColumnType).Ex(context.Background()).Query()

				if currCT != prevCT {
					exprList = append(exprList, d.ModifyColumn(currC, prevC))
				}
				return
			}
			exprList = append(exprList, d.DropColumn(currC))
			return
		}

		if currC.DeprecatedActs == nil {
			exprList = append(exprList, d.AddColumn(currC))
		}
	})

	// indexes
	indexes := map[string]bool{}

	t.Keys.Range(func(key *Key, idx int) {
		name := key.Name
		if key.IsPrimary() {
			name = d.PrimaryKeyName()
		}
		indexes[name] = true

		prevKey := prevT.Key(name)
		if prevKey == nil {
			exprList = append(exprList, d.AddIndex(key))
		} else {
			if !key.IsPrimary() {
				indexDef := key.Def.TableExpr(key.Table).Ex(context.Background()).Query()
				prevIndexDef := prevKey.Def.TableExpr(prevKey.Table).Ex(context.Background()).Query()

				if !strings.EqualFold(indexDef, prevIndexDef) {
					exprList = append(exprList, d.DropIndex(key))
					exprList = append(exprList, d.AddIndex(key))
				}
			}
		}
	})

	prevT.Keys.Range(func(key *Key, idx int) {
		if _, ok := indexes[strings.ToLower(key.Name)]; !ok {
			exprList = append(exprList, d.DropIndex(key))
		}
	})

	return
}

func (t *Table) ColumnsAndValuesByFieldValues(fvs FieldValues) (*Columns, []interface{}) {
	fields := make([]string, 0)
	for name, _ := range fvs {
		fields = append(fields, name)
	}

	sort.Strings(fields)

	cols := &Columns{}
	args := make([]interface{}, 0, len(fvs))

	for _, fieldName := range fields {
		if col := t.ColByFieldName(fieldName); col != nil {
			cols.Add(col)
			args = append(args, fvs[fieldName])
		}
	}
	return cols, args
}

func (t *Table) AssignmentsByFieldValues(fvs FieldValues) Assignments {
	var assignments Assignments
	for name, value := range fvs {
		col := t.ColByFieldName(name)
		if col != nil {
			assignments = append(assignments, col.ValueBy(value))
		}
	}
	return assignments
}

type Tables struct {
	lst    *list.List
	tables map[string]*list.Element
	models map[string]*list.Element
}

func (t *Tables) TableNames() []string {
	if t.lst == nil {
		return nil
	}
	names := make([]string, 0, t.lst.Len())
	t.Range(func(tbl *Table, _ int) {
		names = append(names, tbl.Name)
	})
	return names
}

func (t *Tables) Add(tables ...*Table) {
	if t.lst == nil {
		t.lst = list.New()
		t.tables = make(map[string]*list.Element)
		t.models = make(map[string]*list.Element)
	}
	for _, tbl := range tables {
		if tbl == nil {
			continue
		}
		if _, ok := t.tables[tbl.Name]; ok {
			t.Remove(tbl.Name)
		}
		e := t.lst.PushBack(tbl)
		t.tables[tbl.Name] = e
		if tbl.ModelName != "" {
			t.models[tbl.ModelName] = e
		}
	}
}

func (t *Tables) Table(name string) *Table {
	if t.tables != nil {
		if c, ok := t.tables[name]; ok {
			return c.Value.(*Table)
		}
	}
	return nil
}

func (t *Tables) Model(typename string) *Table {
	if t.models != nil {
		if c, ok := t.models[typename]; ok {
			return c.Value.(*Table)
		}
	}
	return nil
}

func (t *Tables) Remove(name string) {
	if t.tables != nil {
		if e, exists := t.tables[name]; exists {
			t.lst.Remove(e)
			delete(t.tables, name)
		}
	}
}

func (t *Tables) Range(f func(*Table, int)) {
	if t.lst != nil {
		i := 0
		for e := t.lst.Front(); e != nil; e = e.Next() {
			f(e.Value.(*Table), i)
			i++
		}
	}
}
