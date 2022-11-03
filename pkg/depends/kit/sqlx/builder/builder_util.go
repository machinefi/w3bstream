package builder

import (
	"context"
	"go/ast"
	"reflect"
	"strings"
	"sync"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/clone"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type FieldValues map[string]interface{}

type ctxKeyTableName struct{}

func WithTableName(tbl string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, ctxKeyTableName{}, tbl)
	}
}

func TableNameFromContext(ctx context.Context) (string, bool) {
	tbl, ok := ctx.Value(ctxKeyTableName{}).(string)
	return tbl, ok
}

type ctxKeyTableAlias struct{}

func WithTabelAlias(tbl string) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, ctxKeyTableAlias{}, tbl)
	}
}

func TableAliasFromContext(ctx context.Context) (string, bool) {
	tbl, ok := ctx.Value(ctxKeyTableAlias{}).(string)
	return tbl, ok
}

func ColumnsByStruct(v interface{}) *Ex {
	ctx := context.Background()
	fields := StructFieldsFor(ctx, typesx.FromReflectType(reflect.TypeOf(v)))

	e := Expr("")
	e.Grow(len(fields))

	i := 0
	ForEachFieldValue(ctx, reflect.ValueOf(v), func(f *FieldValue) {
		if i > 0 {
			e.WriteQuery(", ")
		}
		if f.TableName != "" {
			e.WriteQuery(f.TableName)
			e.WriteQueryByte('.')
			e.WriteQuery(f.Field.Name)
			e.WriteQuery(" AS ")
			e.WriteQuery(f.TableName)
			e.WriteQuery("__")
			e.WriteQuery(f.Field.Name)
		} else {
			e.WriteQuery(f.Field.Name)
		}
		i++
	})

	return e
}

func ForEachFieldValue(ctx context.Context, v interface{}, fn func(*FieldValue)) {
	rv, ok := v.(reflect.Value)
	if ok {
		if rv.Kind() == reflect.Ptr && rv.IsNil() {
			rv.Set(reflectx.New(rv.Type()))
		}
		v = rv.Interface()
	}
	if m, ok := v.(Model); ok {
		ctx = WithTableName(m.TableName())(ctx)
	}

	fields := StructFieldsFor(ctx, typesx.FromReflectType(reflect.TypeOf(v)))
	rv = reflectx.Indirect(reflect.ValueOf(v))

	for i := range fields {
		f := fields[i]
		tag := f.Tags["db"]

		if tag.HasFlag("deprecated") {
			continue
		}

		if alias, ok := f.Tags["alias"]; ok {
			ctx = WithTabelAlias(alias.Name())(ctx)
		} else {
			if len(f.ModelLoc) > 0 {
				if fmv := f.FieldModelValue(rv); fmv.IsValid() {
					if m, ok := fmv.Interface().(Model); ok {
						ctx = WithTableName(m.TableName())(ctx)
					}
				}
			}
		}
		sf := &FieldValue{
			Field: *f,
			Value: f.FieldValue(rv),
		}
		sf.TableName, _ = TableNameFromContext(ctx)
		if alias, ok := TableAliasFromContext(ctx); ok && alias != "" {
			sf.TableName = alias
		}
		fn(sf)
	}
}

func GetColumnName(name, tag string) string {
	i := strings.Index(tag, ",")
	if tag != "" {
		if i == -1 {
			return strings.ToLower(tag)
		}
		if i > 0 {
			return strings.ToLower(tag[0:i])
		}
	}
	return "f_" + strings.ToLower(name)
}

func ToMap(lst []string) map[string]bool {
	m := make(map[string]bool)
	for _, name := range lst {
		m[name] = true
	}
	return m
}

func FieldValueFromStruct(v interface{}, names []string) FieldValues {
	fvs := FieldValues{}
	rv := reflect.Indirect(reflect.ValueOf(v))
	m := ToMap(names)
	ForEachFieldValue(context.Background(), rv, func(fv *FieldValue) {
		if name := fv.Field.FieldName; m != nil && m[name] {
			fvs[name] = fv.Value.Interface()
		}
	})
	return fvs
}

func FieldValueFromStructByNoneZero(v interface{}, excludes ...string) FieldValues {
	fvs := FieldValues{}
	rv := reflect.Indirect(reflect.ValueOf(v))
	m := ToMap(excludes)
	ForEachFieldValue(context.Background(), rv, func(fv *FieldValue) {
		name := fv.Field.FieldName
		if !reflectx.IsEmptyValue(fv.Value) || m != nil && m[name] {
			fvs[name] = fv.Value.Interface()
		}
	})
	return fvs
}

func TableFromModel(m Model) *Table {
	t := reflect.TypeOf(m)
	if t.Kind() != reflect.Ptr {
		panic("model must be a ptr")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		panic("model must be a struct")
	}
	tbl := T(m.TableName())
	tbl.Model = m
	ScanDefToTable(tbl, m)
	return tbl
}

func ScanDefToTable(tbl *Table, i interface{}) {
	t := typesx.DeRef(typesx.FromReflectType(reflect.TypeOf(i)))
	EachField(context.Background(), t,
		func(f *StructField) bool {
			tbl.AddCol(&Column{
				FieldName:  f.FieldName,
				Name:       f.Name,
				ColumnType: &f.ColumnType,
			})
			return true
		},
	)

	if with, ok := i.(WithTableDesc); ok {
		tbl.Desc = with.TableDesc()
	}
	if with, ok := i.(WithComments); ok {
		for name, comment := range with.Comments() {
			if col := tbl.ColByFieldName(name); col != nil {
				col.Comment = comment
			}
		}
	}
	if with, ok := i.(WithColDesc); ok {
		for name, desc := range with.ColDesc() {
			if col := tbl.ColByFieldName(name); col != nil {
				col.Desc = desc
			}
		}
	}
	if with, ok := i.(WithColRel); ok {
		for name, rel := range with.ColRel() {
			if col := tbl.ColByFieldName(name); col != nil {
				col.Rel = rel
			}
		}
	}
	if with, ok := i.(WithPrimaryKey); ok {
		tbl.AddKey(&Key{
			Name:     "primary",
			IsUnique: true,
			Def:      *ParseIndexDef(with.PrimaryKey()...),
		})
	}
	if with, ok := i.(WithUniqueIndexes); ok {
		for _index, names := range with.UniqueIndexes() {
			name, method := SplitIndexNameAndMethod(_index)
			tbl.AddKey(&Key{
				Name:     name,
				Method:   method,
				IsUnique: true,
				Def:      *ParseIndexDef(names...),
			})
		}
	}
	if with, ok := i.(WithIndexes); ok {
		for _index, names := range with.Indexes() {
			name, method := SplitIndexNameAndMethod(_index)
			tbl.AddKey(&Key{
				Name:   name,
				Method: method,
				Def:    *ParseIndexDef(names...),
			})
		}
	}
}

// SplitIndexNameAndMethod @def index name/method
func SplitIndexNameAndMethod(v string) (string, string) {
	parts := strings.Split(v, "/")
	name := strings.ToLower(parts[0])
	method := ""
	if len(parts) > 1 {
		method = parts[1]
	}
	return name, method
}

type IndexDefine struct {
	Kind   string
	Name   string
	Method string
	IndexDef
}

func (i IndexDefine) ID() string {
	if i.Method != "" {
		return i.Name + "/" + i.Method
	}
	return i.Name
}

func ParseIndexDefine(def string) *IndexDefine {
	d := &IndexDefine{}
	for i := strings.Index(def, " "); i != -1; i = strings.Index(def, " ") {
		part := def[0:i]
		if part != "" {
			if d.Kind == "" {
				d.Kind = part
			} else if d.Name == "" && d.Kind != "primary" {
				d.Name, d.Method = SplitIndexNameAndMethod(part)
			} else {
				break
			}
		}
		def = def[i+1:]
	}
	d.IndexDef = *ParseIndexDef(strings.TrimSpace(def))
	return d
}

type FieldValue struct {
	Field     StructField
	TableName string
	Value     reflect.Value
}

type FieldsFactory struct {
	cache sync.Map
}

// gFields default global field cache
var gFields = &FieldsFactory{}

func StructFieldsFor(ctx context.Context, t typesx.Type) []*StructField {
	return gFields.TableFieldsFor(ctx, t)
}

func (ft *FieldsFactory) TableFieldsFor(ctx context.Context, t typesx.Type) []*StructField {
	t = typesx.DeRef(t)
	i := t.Unwrap() // underlying

	if v, ok := ft.cache.Load(i); ok {
		return v.([]*StructField)
	}

	sf := make([]*StructField, 0)
	EachField(ctx, t, func(f *StructField) bool {
		name := f.Tags["db"]
		if name != "" && name != "-" {
			sf = append(sf, f)
		}
		return true
	})
	ft.cache.Store(i, sf)
	return sf
}

func EachField(ctx context.Context, t typesx.Type, each func(*StructField) bool) {
	must.BeTrue(t.Kind() == reflect.Struct)

	var walk func(t typesx.Type, modelLoc []int, parents ...int)

	walk = func(t typesx.Type, modelLoc []int, parents ...int) {
		if t.Implements(typesx.FromReflectType(RtModel)) {
			modelLoc = parents
		}
		for i := 0; i < t.NumField(); i++ {
			fi := t.Field(i)
			if !ast.IsExported(fi.Name()) {
				continue
			}

			loc := append(parents, i)
			tags := reflectx.ParseStructTag(string(fi.Tag()))
			name := fi.Name()
			tag, has := tags["db"]
			if has {
				if tagName := tag.Name(); tagName == "-" {
					continue
				} else {
					if tagName != "" {
						name = tagName
					}
				}
			}

			if !has && (fi.Anonymous() || fi.Type().Name() == fi.Name()) {
				ft := fi.Type()
				if !ft.Implements(typesx.FromReflectType(RtDriverValuer)) {
					for ft.Kind() == reflect.Ptr {
						ft = ft.Elem()
					}
					if ft.Kind() == reflect.Struct {
						walk(ft, modelLoc, loc...)
						continue
					}
				}
			}

			p := &StructField{
				Name:       strings.ToLower(name),
				FieldName:  fi.Name(),
				Type:       fi.Type(),
				Field:      fi,
				Tags:       tags,
				Loc:        clone.Ints(loc),
				ModelLoc:   clone.Ints(modelLoc),
				ColumnType: *AnalyzeColumnType(fi.Type(), string(tag)),
			}
			if !each(p) {
				break
			}
		}
	}

	walk(t, []int{})
}

type StructField struct {
	Name       string
	FieldName  string
	Type       typesx.Type
	Field      typesx.StructField
	Tags       map[string]reflectx.StructTag
	Loc        []int
	ModelLoc   []int
	ColumnType ColumnType
}

func fieldValue(v reflect.Value, locs []int) reflect.Value {
	n := len(locs)
	if n == 0 {
		return v
	}
	if n < 0 {
		return reflect.Value{}
	}
	v = reflectx.Indirect(v)
	fv := v
	for i := 0; i < n; i++ {
		loc := locs[i]
		fv = fv.Field(loc)
		if i < n-1 {
			for fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					fv.Set(reflectx.New(fv.Type()))
				}
				fv = fv.Elem()
			}
		}
	}
	return fv
}

func (sf *StructField) FieldValue(v reflect.Value) reflect.Value {
	return fieldValue(v, sf.Loc)
}

func (sf *StructField) FieldModelValue(v reflect.Value) reflect.Value {
	return fieldValue(v, sf.ModelLoc)
}
