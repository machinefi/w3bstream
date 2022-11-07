package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/saitofun/qlib/util/qnaming"

	g "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/kit/modelgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

type Table struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Defs    []Def  `json:"defs"`
	Cols    []Col  `json:"cols"`
	cols    map[string]int
}

func (t *Table) TableName() string {
	return "t_" + strings.ToLower(strings.TrimPrefix(t.Name, "t_"))
}

func (t *Table) TypeName() string {
	return qnaming.UpperCamelCase(strings.TrimPrefix(t.Name, "t_"))
}

type Def struct {
	Type DefType  `json:"type"`
	Name string   `json:"name"`
	Cols []string `json:"cols"`
}

type Col struct {
	Name        string      `json:"name"`
	Type        Datatype    `json:"type"`
	Constraints []Constrain `json:"constraints"`
	Comment     string      `json:"comment"`
}

func (c *Col) FieldName() string {
	return qnaming.UpperCamelCase(strings.Trim(c.Name, "f_"))
}

func (c *Col) ColumnName() string {
	return "f_" + qnaming.LowerSnakeCase(strings.TrimPrefix(c.Name, "f_"))
}

func (c *Col) Tag() string {
	constrains := []string{c.ColumnName()}
	for _, c := range c.Constraints {
		if tag := c.Tag(); tag != "" {
			constrains = append(constrains, tag)
		}
	}
	return `db:"` + strings.Join(constrains, ",") + `"`
}

type Constrain struct {
	Type  ConstrainType `json:"type"`
	Value interface{}   `json:"value"`
}

func (c *Constrain) Tag() string {
	return c.Type.Tag("", c.Value)
}

type Schema struct {
	Project string
	Version string
	Tables  []Table
	files   map[string]*g.File
}

func Load(content []byte, root, project, version string) (*Schema, error) {
	buf := bytes.NewBuffer(content)

	c := &Schema{
		Project: project,
		Version: version,
		files:   make(map[string]*g.File),
	}
	err := json.NewDecoder(buf).Decode(&c.Tables)
	if err != nil {
		return nil, err
	}
	for i := range c.Tables {
		c.Tables[i].cols = make(map[string]int)
		for j := range c.Tables[i].Cols {
			c.Tables[i].cols[c.Tables[i].Cols[j].Name] = j
		}
	}

	pkgPath := filepath.Join(root, "models", c.Project)
	if c.Version != "" {
		pkgPath = filepath.Join(pkgPath, c.Version)
	} else {
		pkgPath = filepath.Join(pkgPath, "0.0.0")
	}

	if err := c.WriteSchema(pkgPath); err != nil {
		return nil, err
	}

	if err := c.WriteTables(pkgPath); err != nil {
		return nil, err
	}

	return c, nil
}

func LoadConfigFrom(path, root, project, version string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return Load(content, root, project, version)
}

func (c *Schema) SchemaName() string {
	if c.Version == "" {
		return c.Project + "_" + "0.0.0"
	}
	return c.Project + "_" + c.Version
}

func (c *Schema) File(name string) *g.File {
	if f, ok := c.files[name]; !ok {
		panic("not define table " + name)
	} else {
		return f
	}
}

func (c *Schema) Comments() map[string][]string {
	comments := make(map[string][]string)
	for _, t := range c.Tables {
		types := make([]string, 0)
		maxLenDefs := 0
		names := make([]string, 0)
		maxLenName := 0
		cols := make([][]string, 0)

		for _, d := range t.Defs {
			def := strings.ToLower(d.Type.String())
			types = append(types, def)
			if len(def) > maxLenDefs {
				maxLenDefs = len(def)
			}
			names = append(names, d.Name)
			if len(d.Name) > maxLenName {
				maxLenName = len(d.Name)
			}
			fields := make([]string, 0)
			for _, col := range d.Cols {
				fields = append(fields, t.Cols[t.cols[col]].FieldName())
			}
			cols = append(cols, fields)
		}

		for i := 0; i < len(types); i++ {
			types[i] += strings.Repeat(" ", maxLenDefs-len(types[i]))
			names[i] += strings.Repeat(" ", maxLenName-len(names[i]))
		}

		defs := make([]string, 0)
		if t.Comment != "" {
			defs = append(defs, t.TypeName()+" "+t.Comment)
		}
		for i := 0; i < len(types); i++ {
			defs = append(defs,
				fmt.Sprintf("@def %s %s %s",
					types[i], names[i], strings.Join(cols[i], " "),
				),
			)
		}
		comments[t.Name] = defs

	}
	return comments
}

func (c *Schema) SnippetStruct(comments map[string][]string) map[string]g.Snippet {
	snippets := make(map[string]g.Snippet)
	for _, t := range c.Tables {
		fields := make([]*g.SnippetField, 0)
		for _, f := range t.Cols {
			sf := (&g.SnippetField{
				Type:  f.Type.CodeGenType(c.File(t.Name)),
				Names: []*g.SnippetIdent{g.Ident(f.FieldName())},
			}).WithTag(f.Tag())
			if f.Comment != "" {
				sf = sf.WithOneLineComment(f.Comment)
			}
			fields = append(fields, sf)
		}

		decl := g.DeclType(g.Var(g.Struct(fields...), t.TypeName()))
		if cmts, ok := comments[t.Name]; ok {
			decl = decl.WithComments(cmts...)
		}
		snippets[t.Name] = decl
	}
	return snippets
}

func (c *Schema) SnippetSchema(f *g.File) g.Snippet {
	return g.DeclVar(g.Assign(g.Ident("DB")).By(
		g.Exprer(
			f.Use("github.com/machinefi/w3bstream/pkg/models", "DB")+`.?(?)`,
			g.Ident("WithSchema"),
			g.Valuer(c.SchemaName())),
	))
}

func (c *Schema) WriteSchema(root string) error {
	f := g.NewFile(c.Project, filepath.Join(root, "schema.go"))
	f.WriteSnippet(c.SnippetSchema(f))
	_, err := f.Write()
	return err
}

func (c *Schema) WriteTables(root string) error {
	for i := range c.Tables {
		c.files[c.Tables[i].Name] = g.NewFile(
			c.Project,
			filepath.Join(
				root,
				qnaming.LowerSnakeCase(c.Tables[i].Name)+".go",
			),
		)
	}

	snippets := c.SnippetStruct(c.Comments())

	for name, file := range c.files {
		file.WriteSnippet(snippets[name])
		if _, err := file.Write(); err != nil {
			panic(err)
		}
	}

	cwd, _ := os.Getwd()

	for _, t := range c.Tables {
		pkg, err := pkgx.LoadFrom(root)
		if err != nil {
			return err
		}
		mg := modelgen.New(pkg)
		mg.WithComments = true
		mg.WithTableName = true
		mg.WithTableInterfaces = true
		mg.StructName = t.TypeName()
		mg.Database = "DB"
		mg.Scan()
		mg.Output(cwd)
	}
	return nil
}
