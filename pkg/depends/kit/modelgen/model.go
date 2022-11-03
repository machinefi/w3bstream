package modelgen

import (
	"bytes"
	"fmt"
	"go/types"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/mapx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

func NewModel(pkg *pkgx.Pkg, tn *types.TypeName, doc string, cfg *Config) *Model {
	cfg.SetDefault()
	m := Model{
		TypeName: tn,
		Config:   cfg,
		Table:    builder.T(cfg.TableName),
	}

	// parse comments for each struct field
	{

		t := tn.Type().Underlying().Underlying().(*types.Struct) // struct type
		p := pkg.PkgByPath(tn.Pkg().Path())                      // package info
		// for each struct field var
		each := func(v *types.Var, name string, tag string) {
			col := builder.Col(name).Field(v.Name()).Type("", tag)
			for ident, obj := range p.TypesInfo.Defs {
				if obj != v {
					continue
				}
				rel, lines := parseColRelFromDoc(pkg.CommentsOf(ident))
				if rel != "" {
					if path := strings.Split(rel, "."); len(path) == 2 {
						col.Rel = path
					} else {
						continue
					}
				}
				if len(lines) > 0 {
					col.Comment, col.Desc = lines[0], lines
				}
			}
			m.AddColumnAndField(col, v)
		}
		forEachStructField(t, each)
	}

	m.HasDeletedAt = m.Table.ColByFieldName(m.FieldKeyDeletedAt) != nil
	m.HasCreatedAt = m.Table.ColByFieldName(m.FieldKeyCreatedAt) != nil
	m.HasUpdatedAt = m.Table.ColByFieldName(m.FieldKeyUpdatedAt) != nil

	keys, lines := parseKeysFromDoc(doc)
	m.Keys = keys
	if len(lines) > 0 {
		m.Desc = lines
	}

	if m.HasDeletedAt {
		m.Keys.WithSoftDelete(m.FieldKeyDeletedAt)
	}
	m.Keys.Bind(m.Table)

	if col := m.Table.AutoIncrement(); col != nil {
		m.HasAutoIncrement = true
		m.FieldKeyAutoIncrement = col.FieldName
	}

	return &m
}

type Model struct {
	*types.TypeName
	*Config
	*Keys
	*builder.Table
	Fields                map[string]*types.Var
	FieldKeyAutoIncrement string
	HasDeletedAt          bool
	HasCreatedAt          bool
	HasUpdatedAt          bool
	HasAutoIncrement      bool
}

func (m *Model) AddColumnAndField(col *builder.Column, tpe *types.Var) {
	m.Table.Columns.Add(col)
	if m.Fields == nil {
		m.Fields = map[string]*types.Var{}
	}
	m.Fields[col.FieldName] = tpe
}

func (m *Model) GetComments() map[string]string {
	comments := map[string]string{}
	m.Columns.Range(func(c *builder.Column, _ int) {
		if c.Comment != "" {
			comments[c.FieldName] = c.Comment
		}
	})
	return comments
}

func (m *Model) GetColDesc() map[string][]string {
	desc := map[string][]string{}
	m.Columns.Range(func(c *builder.Column, _ int) {
		if len(c.Desc) > 0 {
			desc[c.FieldName] = c.Desc
		}
	})
	return desc
}

func (m *Model) GetColRel() map[string][]string {
	rel := map[string][]string{}
	m.Columns.Range(func(c *builder.Column, _ int) {
		if len(c.Rel) == 2 {
			rel[c.FieldName] = c.Rel
		}
	})
	return rel
}

func (m *Model) GetIndexFieldNames() []string {
	names := make([]string, 0)
	m.Table.Keys.Range(func(k *builder.Key, _ int) {
		names = append(names, k.Def.FieldNames...)
	})
	names = uniqueStrings(names)
	names = filterStrings(names, func(s string, i int) bool {
		return m.HasDeletedAt && s != m.FieldKeyDeletedAt || !m.HasDeletedAt
	})
	sort.Strings(names)
	return names
}

func (m *Model) Type() g.SnippetType {
	return g.Type(m.StructName)
}

func (m *Model) IteratorType() g.SnippetType {
	return g.Type(m.StructName + "Iterator")
}

func (m *Model) IteratorPtrType() g.SnippetType {
	return g.Star(m.IteratorType())
}

func (m *Model) PtrType() g.SnippetType {
	return g.Star(m.Type())
}

func (m *Model) VarTable() string {
	return m.StructName + "Table"
}

func (m *Model) FileType(f *g.File, fn string) g.SnippetType {
	field, ok := m.Fields[fn]
	if !ok {
		return nil
	}
	typ := field.Type().String()
	if strings.Contains(typ, ".") {
		pkg, name := pkgx.ImportPathAndExpose(typ)
		if pkg != m.TypeName.Pkg().Path() {
			return g.Type(f.Use(pkg, name))
		}
		return g.Type(name)
	}
	return g.BuiltInType(typ)
}

// Snippets:

// SnippetTableInstanceAndInit generated below
// `Model`Table *builder.Table
// func init() { `Model`Table = `DB`.Register(&`Model`{})
func (m *Model) SnippetTableInstanceAndInit(f *g.File) []g.Snippet {
	return []g.Snippet{
		g.DeclVar(
			g.Var(
				g.Star(g.Type(f.Use(BuilderPkg, "Table"))),
				m.VarTable(),
			),
		),

		g.Func().Named("init").
			Do(
				g.Exprer(
					`?=?.Register(&?{})`,
					g.Ident(m.VarTable()),
					g.Ident(m.Database),
					g.Ident(m.StructName),
				),
			),
	}

}

// SnippetTableIteratorAndMethods generate below
// `Model`Iterator struct{}
// func (`Model`Table) New() interface{} { return &`Model`{} }
// func (`Model`Table) Resolve(v interface{}) *`Model` { return v.(*`Model`) }
func (m *Model) SnippetTableIteratorAndMethods(_ *g.File) []g.Snippet {
	return []g.Snippet{
		g.DeclType(g.Var(g.Struct(), string(m.IteratorType().Bytes()))),

		g.Func().
			Named("New").
			MethodOf(g.Var(m.IteratorPtrType())).
			Return(g.Var(g.Interface())).
			Do(g.Return(g.Exprer("&?{}", m.Type()))),

		g.Func(g.Var(g.Interface(), "v")).
			Named("Resolve").
			MethodOf(g.Var(m.IteratorPtrType())).
			Return(g.Var(m.PtrType())).
			Do(g.Return(g.Exprer("v.(?)", m.PtrType()))),
	}
}

// SnippetTableName generate below
// TableName implements builder.Model
// func (`Model`) TableName() string
func (m *Model) SnippetTableName(f *g.File) g.Snippet {
	return g.Func().Named("TableName").MethodOf(g.Var(m.PtrType())).
		Return(g.Var(g.String)).
		Do(g.Return(f.Value(m.Config.TableName)))
}

// SnippetTableDesc generate below
// TableDesc implements builder.WithTableDesc
// func (`Model`) TableDesc() []string
func (m *Model) SnippetTableDesc(f *g.File) g.Snippet {
	if len(m.Table.Desc) == 0 {
		return nil
	}
	return g.Func().Named("TableDesc").MethodOf(g.Var(m.PtrType())).
		Return(g.Var(g.Slice(g.String))).
		Do(g.Return(f.Value(m.Table.Desc)))
}

// SnippetComments generate below
// Comments implements builder.WithComments
// func(`Model`) Comments() map[string]string    // field_name: comment line
func (m *Model) SnippetComments(f *g.File) g.Snippet {
	if !m.WithComments {
		return nil
	}
	return g.Func().Named("Comments").MethodOf(g.Var(m.PtrType())).
		Return(g.Var(g.Map(g.String, g.String))).
		Do(g.Return(f.Value(m.GetComments())))
}

// SnippetColDesc generate below
// ColDesc implements builder.WithColDesc
// func (`Model`) ColDesc() map[string][]string
func (m *Model) SnippetColDesc(f *g.File) g.Snippet {
	desc := m.GetColDesc()
	return g.Func().Named("ColDesc").MethodOf(g.Var(m.PtrType())).
		Return(g.Var(g.Map(g.String, g.Slice(g.String)))).
		Do(g.Return(f.Value(desc)))
}

// SnippetColRel generate below
// ColRel implements builder.WithColRel
// func (`Model`) ColRel() map[string][]string
func (m *Model) SnippetColRel(f *g.File) g.Snippet {
	rel := m.GetColRel()
	return g.Func().Named("ColRel").MethodOf(g.Var(m.PtrType())).
		Return(g.Var(g.Map(g.String, g.Slice(g.String)))).
		Do(g.Return(f.Value(rel)))
}

// SnippetPrimaryKey generate below
// PrimaryKey implements builder.WithPrimaryKey
// func(`Model`) PrimaryKey() []string
func (m *Model) SnippetPrimaryKey(f *g.File) g.Snippet {
	if len(m.Keys.Primary) == 0 {
		return nil
	}
	return g.Func().Named("PrimaryKey").MethodOf(g.Var(m.PtrType())).
		Return(g.Var(g.Slice(g.String))).
		Do(g.Return(f.Value(m.Keys.Primary)))
}

// SnippetIndexes generate below
// Indexes implements builder.WithIndexes
// func(`Model`) Indexes() builder.Indexes
func (m *Model) SnippetIndexes(f *g.File) g.Snippet {
	if len(m.Keys.Indexes) == 0 {
		return nil
	}
	return g.Func().Named("Indexes").MethodOf(g.Var(m.PtrType())).
		Return(g.Var(g.Type(f.Use(BuilderPkg, "Indexes")))).
		Do(g.Return(f.Value(m.Keys.Indexes)))
}

// SnippetIndexFieldNames generate below
// func (m *`Model`) IndexFieldNames() []string  // index field name list
func (m *Model) SnippetIndexFieldNames(f *g.File) g.Snippet {
	return g.Func().Named("IndexFieldNames").MethodOf(g.Var(m.PtrType(), "m")).
		Return(g.Var(g.Slice(g.String))).
		Do(g.Return(f.Value(m.GetIndexFieldNames())))
}

// SnippetUniqueIndexes generate below
// UniqueIndexes() implements `builder.WithUniqueIndexes`
// func(`Model`) UniqueIndexes() builder.Indexes [db_index_name->field_names]
// func(`Model`) UniqueIndexXXX() string;        [for each unique index, XXX is index name]
func (m *Model) SnippetUniqueIndexes(f *g.File) []g.Snippet {
	if len(m.UniqueIndexes) == 0 {
		return nil
	}
	snippets := make([]g.Snippet, 0, len(m.UniqueIndexes)+1)
	names := make([]string, 0, len(m.UniqueIndexes))
	for name := range m.Keys.UniqueIndexes {
		names = append(names, name)
	}
	sort.Strings(names)
	snippets = append(snippets,
		g.Func().Named("UniqueIndexes").MethodOf(g.Var(m.PtrType())).
			Return(g.Var(g.Type(f.Use(BuilderPkg, "Indexes")))).
			Do(g.Return(f.Value(m.Keys.UniqueIndexes))),
	)
	for _, name := range names {
		fn := "UniqueIndex" + stringsx.UpperCamelCase(name)
		snippets = append(snippets,
			g.Func().Named(fn).MethodOf(g.Var(m.PtrType())).
				Return(g.Var(g.String)).
				Do(g.Return(f.Value(name))),
		)
	}
	return snippets
}

// SnippetFieldMethods generate below
// func (m *`Model`) ColXXX() *builder.Column  // FOR EACH Field, return field column
// func (m `Model`)  FieldXXX() string         // FOR EACH Field, return field name
func (m *Model) SnippetFieldMethods(f *g.File) []g.Snippet {
	snippets := make([]g.Snippet, 0, 2*m.Columns.Len())
	m.Columns.Range(func(c *builder.Column, _ int) {
		if c.DeprecatedActs != nil {
			return // TODO generate deprecated actions
		}
		fn := "Col" + c.FieldName
		snippets = append(snippets,
			g.Func().Named(fn).MethodOf(g.Var(m.PtrType(), "m")).
				Return(g.Var(g.Star(g.Type(f.Use(BuilderPkg, "Column"))))).
				Do(
					g.Return(
						g.Exprer(
							"?.ColByFieldName(m.Field"+c.FieldName+"())",
							g.Ident(m.VarTable()),
						),
					),
				),
		)
		fn = "Field" + c.FieldName
		snippets = append(snippets,
			g.Func().Named(fn).MethodOf(g.Var(m.PtrType())).
				Return(g.Var(g.String)).
				Do(g.Return(f.Value(c.FieldName))),
		)
	})
	return snippets
}

// SnippetCondByValue generate below
// func (m *`Model`) CondByValue(DBExecutor) builder.SqlCondition // condition of this
func (m *Model) SnippetCondByValue(f *g.File) g.Snippet {
	if !m.WithMethods {
		return nil
	}

	return g.Func(g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`)).
		Named("CondByValue").MethodOf(g.Var(m.PtrType(), `m`)).
		Return(g.Var(g.Type(f.Use(BuilderPkg, `SqlCondition`)))).
		Do(
			g.DeclVar(
				g.Assign(g.Var(nil, `tbl`)).By(g.Ref(g.Ident(`db`), g.Call(`T`, g.Ident(`m`)))),
				g.Assign(g.Var(nil, `fvs`)).By(g.Call(f.Use(BuilderPkg, `FieldValueFromStructByNoneZero`), g.Ident(`m`))),
				m.DeletedAtCondInitial(f, g.Ident(`tbl`), g.Ident(`cond`)),
			),
			g.Exprer(`
for _, fn := range m.IndexFieldNames() {
if v, ok := fvs[fn]; ok {
cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
delete(fvs, fn)
}
}
if len(cond) == 0 {
panic(`+f.Use(`fmt`, "Errorf")+`("no field for indexes has value"))
}
for fn, v := range fvs {
cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
}`,
			),
			g.Return(g.Call(f.Use(BuilderPkg, "And"), g.Exprer(`cond...`))),
		)
}

func (m *Model) SetCreatedSnippet(f *g.File) g.Snippet {
	if !m.HasCreatedAt {
		return nil
	}
	return g.Exprer(`
if m.` + m.FieldKeyCreatedAt + `.IsZero() {
m.` + m.FieldKeyCreatedAt + `.Set(` + f.Use(`time`, `Now`) + `())
}`)
}

func (m *Model) SetUpdatedSnippet(f *g.File) g.Snippet {
	if !m.HasUpdatedAt {
		return nil
	}
	return g.Exprer(`
if m.` + m.FieldKeyUpdatedAt + `.IsZero() {
m.` + m.FieldKeyUpdatedAt + `.Set(` + f.Use(`time`, `Now`) + `())
}`)
}

func (m *Model) SetUpdatedSnippetForFVs(f *g.File, fvs *g.SnippetIdent) g.Snippet {
	if !m.HasUpdatedAt {
		return nil
	}
	return f.Expr(`
if _, ok := ?[?]; !ok {
?[?] = ?{Time: ?()}
}`,
		fvs, g.Valuer(m.FieldKeyUpdatedAt),
		fvs, g.Valuer(m.FieldKeyUpdatedAt),
		m.FileType(f, m.FieldKeyUpdatedAt), g.Ident(f.Use(`time`, `Now`)),
	)
}

func (m *Model) DeletedAtCondInitial(f *g.File, tbl, cond *g.SnippetIdent) g.SnippetSpec {
	if !m.HasDeletedAt {
		return g.Assign(cond).
			By(g.Call(`make`, g.Slice(g.Type(f.Use(BuilderPkg, `SqlCondition`))), f.Value(0)))
	}
	// cond = []builder.SqlCondition{tbl.ColByFieldName("DeletedAt").Eq(0)}
	return g.Assign(cond).
		By(g.Exprer(
			`[]`+f.Use(BuilderPkg, `SqlCondition`)+`{?.ColByFieldName(?).Eq(0)}`,
			tbl, f.Value(m.FieldKeyDeletedAt),
		))
}

func (m *Model) DeletedAtCondAttach(f *g.File, tbl, cond *g.SnippetIdent) g.SnippetSpec {
	if !m.HasDeletedAt {
		return nil
	}
	return g.Assign(cond).By(g.Call(
		f.Use(BuilderPkg, `And`),
		g.Exprer(`?.ColByFieldName(?).Eq(0)`, tbl, m.FieldKeyDeletedAt),
		cond,
	))
}

func (m *Model) SetDeletedSnippetForFVs(f *g.File, fvs *g.SnippetIdent) g.Snippet {
	if !m.HasUpdatedAt {
		return nil
	}
	return g.Exprer(`
if _, ok := ?[?]; !ok {
?[?] = ?{Time: ?()}
}`,
		fvs, g.Valuer(m.FieldKeyDeletedAt),
		fvs, g.Valuer(m.FieldKeyDeletedAt),
		m.FileType(f, m.FieldKeyDeletedAt), g.Ident(f.Use(`time`, `Now`)),
	)
}

// SnippetCreate generate below
// Create to create record by value m
// func (m *`Model`) Create(DBExecutor) error // Create by this
func (m *Model) SnippetCreate(f *g.File) g.Snippet {
	if !m.WithMethods {
		return nil
	}
	return g.Func(g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`)).Named(`Create`).
		MethodOf(g.Var(m.PtrType(), `m`)).
		Return(g.Var(g.Error)).
		Do(
			m.SetCreatedSnippet(f),
			m.SetUpdatedSnippet(f),
			g.Exprer(`
_, err := db.Exec(?(db, m, nil))
return err`, g.Ident(f.Use(SQLxPkg, `InsertToDB`))),
		)
}

// SnippetList generate below
// List by condition and additions(offset, size)
// func (m *`Model`) List(DBExecutor, SqlCondition, Additions) []`Model`
func (m *Model) SnippetList(f *g.File) g.Snippet {
	if !m.WithMethods {
		return nil
	}
	return g.Func(
		g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`),
		g.Var(g.Type(f.Use(BuilderPkg, `SqlCondition`)), `cond`),
		g.Var(g.Ellipsis(g.Type(f.Use(BuilderPkg, `Addition`))), `adds`),
	).Named(`List`).MethodOf(g.Var(m.PtrType(), `m`)).
		Return(
			g.Var(g.Slice(g.Type(m.StructName))),
			g.Var(g.Error),
		).
		Do(
			g.DeclVar(
				g.Assign(g.Var(nil, `tbl`)).By(g.Ref(g.Ident(`db`), g.Call(`T`, g.Ident(`m`)))),
				g.Assign(g.Var(nil, `lst`)).By(g.Call(`make`, g.Slice(g.Type(m.StructName)), g.Valuer(0))),
			),
			m.DeletedAtCondAttach(f, g.Ident(`tbl`), g.Ident(`cond`)),
			g.Assign(g.Var(nil, `adds`)).By(g.Exprer(
				`append([]`+f.Use(BuilderPkg, `Addition`)+
					`{`+f.Use(BuilderPkg, `Where`)+`(cond), `+
					f.Use(BuilderPkg, `Comment`)+`("`+m.StructName+`.List")}, adds...)`,
			)),
			g.Exprer(`err := db.QueryAndScan(`+f.Use(BuilderPkg, `Select`)+`(nil).From(tbl, adds...), &lst)`),
			g.Return(g.Ident(`lst`), g.Ident(`err`)),
		)
}

// SnippetCount generate below
// Count by condition and addition(offset, size)
// func (m *`Model`) Count(DBExecutor, SqlCondition, Additions) (int64, error)
func (m *Model) SnippetCount(f *g.File) g.Snippet {
	if !m.WithMethods {
		return nil
	}
	return g.Func(
		g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`),
		g.Var(g.Type(f.Use(BuilderPkg, `SqlCondition`)), `cond`),
		g.Var(g.Ellipsis(g.Type(f.Use(BuilderPkg, `Addition`))), `adds`),
	).Named(`Count`).MethodOf(g.Var(m.PtrType(), `m`)).
		Return(g.Var(g.Int64, `cnt`), g.Var(g.Error, `err`)).
		Do(
			g.Define(g.Var(nil, `tbl`)).By(g.Ref(g.Ident(`db`), g.Call(`T`, g.Ident(`m`)))),
			m.DeletedAtCondAttach(f, g.Ident(`tbl`), g.Ident(`cond`)),
			g.Assign(g.Var(nil, `adds`)).By(g.Exprer(
				`append([]`+f.Use(BuilderPkg, `Addition`)+
					`{`+f.Use(BuilderPkg, `Where`)+`(cond), `+
					f.Use(BuilderPkg, `Comment`)+`("`+m.StructName+`.List")}, adds...)`,
			)),
			g.Exprer(`err = db.QueryAndScan(`+f.Use(BuilderPkg, `Select`)+`(`+f.Use(BuilderPkg, `Count`)+`()).From(tbl, adds...), &cnt)`),
			g.Return(),
		)
}

func IndexCond(f *g.File, fns ...string) string {
	b := bytes.NewBuffer(nil)
	b.WriteString(f.Use(BuilderPkg, "And(\n"))
	for _, fn := range fns {
		b.WriteString(fmt.Sprintf(`tbl.ColByFieldName("%s").Eq(m.%s),`, fn, fn))
		b.WriteRune('\n')
	}
	b.WriteString("),")
	return b.String()
}

// SnippetCRUDByUniqueKeys generate below
// UpdateByXXX to update record by value m
// XXX is UniqueIndexNames contacted by `And`; function name like `UpdateByNameAndIDAnd...`
// func (m *`Model`) UpdateByXXX(DBExecutor, zeros...)
// FetchByXXX to select record by some unique index
// XXX is UniqueIndexNames contacted by `And`; function name like `FetchByNameAndIDAnd...`
// func (m *`Model`) FetchByXXX(DBExecutor) error      // XXX is UniqueIndexNames
// Delete by value condition
// func (m *`Model`) Delete(DBExecutor) error
// DeleteByXXX to delete record by some unique index
// XXX is UniqueIndexNames contacted by `And`; function name like `DeleteByNameAndIDAnd...`
// func (m *`Model`) DeleteByXXX(DBExecutor) error     // XXX is UniqueIndexNames
// SoftDeleteByXXX to update `DeleteAt` flag as current time
// func (m *`Model`) SoftDeleteByXXX(DBExecutor) error // XXX is UniqueIndexNames
func (m *Model) SnippetCRUDByUniqueKeys(f *g.File, keys ...string) []g.Snippet {
	if !m.WithMethods {
		return nil
	}
	fetches := make([]g.Snippet, 0)
	updates := make([]g.Snippet, 0)
	deletes := []g.Snippet{
		g.Func(g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`)).
			Named("Delete").MethodOf(g.Var(m.PtrType(), `m`)).
			Return(g.Var(g.Error)).
			Do(
				g.Exprer(`_, err := db.Exec(
`+f.Use(BuilderPkg, `Delete`)+`().
From(
db.T(m),
`+f.Use(BuilderPkg, `Where`)+`(m.CondByValue(db)),
`+f.Use(BuilderPkg, `Comment`)+`(?),
),
)`,
					f.Value(m.StructName+".Delete"),
				),
				g.Return(g.Ident(`err`)),
			),
	}

	set := mapx.Set[string]{}
	if len(keys) > 0 {
		set, _ = mapx.ToSet(keys, strings.ToLower)
	}

	names := make([]string, 0, m.Table.Keys.Len())

	m.Table.Keys.Range(func(k *builder.Key, _ int) {
		names = append(names, k.Name)
	})
	sort.Strings(names)

	for _, name := range names {
		k := m.Table.Keys.Key(name)
		if !k.IsUnique {
			continue
		}
		if len(set) != 0 && !set[strings.ToLower(k.Name)] {
			continue
		}
		fns := k.Def.FieldNames
		kns := filterStrings(fns, func(s string, _ int) bool {
			return m.HasDeletedAt && s != m.FieldKeyDeletedAt || !m.HasDeletedAt
		})
		if m.HasDeletedAt && k.IsPrimary() {
			fns = append(fns, m.FieldKeyDeletedAt)
		}
		xxx := strings.Join(kns, "And")
		mthNameFetchBy := "FetchBy" + xxx

		// FetchByXXX
		fetches = append(fetches,
			g.Func(g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`)).
				Named(mthNameFetchBy).MethodOf(g.Var(m.PtrType(), `m`)).
				Return(g.Var(g.Error)).
				Do(
					g.Exprer(`tbl := db.T(m)
err := db.QueryAndScan(
`+f.Use(BuilderPkg, `Select`)+`(nil).
From(
tbl,
`+f.Use(BuilderPkg, `Where`)+`(
`+IndexCond(f, fns...)+`
),
`+f.Use(BuilderPkg, `Comment`)+`(?),
),
m,
)`, f.Value(m.StructName+"."+mthNameFetchBy)),
					g.Return(g.Ident(`err`)),
				),
		)
		// UpdateByXXXWithFVs
		mthNameUpdateByWithFVs := "UpdateBy" + xxx + "WithFVs"
		updates = append(updates,
			g.Func(
				g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`),
				g.Var(g.Type(f.Use(BuilderPkg, `FieldValues`)), `fvs`),
			).
				Named(mthNameUpdateByWithFVs).MethodOf(g.Var(m.PtrType(), `m`)).
				Return(g.Var(g.Error)).
				Do(
					m.SetUpdatedSnippetForFVs(f, g.Ident(`fvs`)),
					g.Exprer(`tbl := db.T(m)
res, err := db.Exec(
`+f.Use(BuilderPkg, `Update`)+`(tbl).
Where(
`+IndexCond(f, fns...)+`
`+f.Use(BuilderPkg, `Comment`)+`(?),
).
Set(tbl.AssignmentsByFieldValues(fvs)...),
)
if err != nil {
return err
}
if affected, _ := res.RowsAffected(); affected == 0 {
return m.`+mthNameFetchBy+`(db)
}
return nil`,
						f.Value(m.StructName+"."+mthNameUpdateByWithFVs),
					),
				),
		)

		// UpdateByXXX
		mthNameUpdateBy := "UpdateBy" + xxx
		updates = append(updates,
			g.Func(
				g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`),
				g.Var(g.Ellipsis(g.String), `zeros`),
			).
				Named(mthNameUpdateBy).MethodOf(g.Var(m.PtrType(), `m`)).
				Return(g.Var(g.Error)).
				Do(
					g.Exprer(`fvs := `+f.Use(BuilderPkg, `FieldValueFromStructByNoneZero`)+`(m, zeros...)
return m.`+mthNameUpdateByWithFVs+`(db, fvs)`,
					),
				),
		)

		// DeleteByXXX
		mthNameDeleteBy := "DeleteBy" + xxx
		deletes = append(deletes,
			g.Func(g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`)).
				Named(mthNameDeleteBy).MethodOf(g.Var(m.PtrType(), `m`)).
				Return(g.Var(g.Error)).
				Do(
					g.Exprer(`tbl := db.T(m)
_, err := db.Exec(
`+f.Use(BuilderPkg, `Delete`)+`().
From(
tbl,
`+f.Use(BuilderPkg, `Where`)+`(
`+IndexCond(f, fns...)+`
),
`+f.Use(BuilderPkg, `Comment`)+`(?),
),
)`,
						f.Value(m.StructName+"."+mthNameDeleteBy),
					),
					g.Return(g.Ident(`err`)),
				),
		)

		if m.HasDeletedAt {
			mthNameSoftDeleteBy := "SoftDeleteBy" + xxx
			deletes = append(deletes,
				g.Func(g.Var(g.Type(f.Use(SQLxPkg, `DBExecutor`)), `db`)).
					Named(mthNameSoftDeleteBy).MethodOf(g.Var(m.PtrType(), `m`)).
					Return(g.Var(g.Error)).
					Do(
						g.Exprer(`tbl := db.T(m)
fvs := `+f.Use(BuilderPkg, `FieldValues`)+`{}`),
						m.SetDeletedSnippetForFVs(f, g.Ident(`fvs`)),
						m.SetUpdatedSnippetForFVs(f, g.Ident(`fvs`)),
						g.Exprer(`_, err := db.Exec(
`+f.Use(BuilderPkg, `Update`)+`(db.T(m)).
Where(
`+IndexCond(f, fns...)+`
`+f.Use(BuilderPkg, `Comment`)+`(?),
).
Set(tbl.AssignmentsByFieldValues(fvs)...),
)
return err`,
							f.Value(m.StructName+"."+mthNameSoftDeleteBy)),
					),
			)
		}
	}
	return append(fetches, append(updates, deletes...)...)
}

func (m *Model) WriteTo(f *g.File) {
	snippets := make([]g.Snippet, 0)
	snippets = append(snippets, m.SnippetTableInstanceAndInit(f)...)
	snippets = append(snippets, m.SnippetTableIteratorAndMethods(f)...)
	snippets = append(snippets, m.SnippetTableName(f))
	snippets = append(snippets, m.SnippetTableDesc(f))
	snippets = append(snippets, m.SnippetComments(f))
	snippets = append(snippets, m.SnippetColDesc(f))
	snippets = append(snippets, m.SnippetColRel(f))
	snippets = append(snippets, m.SnippetPrimaryKey(f))
	snippets = append(snippets, m.SnippetIndexes(f))
	snippets = append(snippets, m.SnippetIndexFieldNames(f))
	snippets = append(snippets, m.SnippetUniqueIndexes(f)...)
	snippets = append(snippets, m.SnippetFieldMethods(f)...)
	snippets = append(snippets, m.SnippetCondByValue(f))
	snippets = append(snippets, m.SnippetCreate(f))
	snippets = append(snippets, m.SnippetList(f))
	snippets = append(snippets, m.SnippetCount(f))
	snippets = append(snippets, m.SnippetCRUDByUniqueKeys(f)...)

	f.WriteSnippet(snippets...)
}

var (
	BuilderPkg = "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	SQLxPkg    = "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
)

func init() {
	_, current, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(current), "../sqlx")
	SQLxPkg = must.String(pkgx.PkgIdByPath(dir))
	dir = filepath.Join(filepath.Dir(current), "../sqlx/builder")
	BuilderPkg = must.String(pkgx.PkgIdByPath(dir))
}
