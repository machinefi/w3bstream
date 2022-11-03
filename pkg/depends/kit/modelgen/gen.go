package modelgen

import (
	"go/types"
	"path"
	"path/filepath"

	"github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

type Generator struct {
	Config
	pkg *pkgx.Pkg
	mod *Model
}

func New(pkg *pkgx.Pkg) *Generator { return &Generator{pkg: pkg} }

func (g *Generator) Scan() {
	for ident, obj := range g.pkg.TypesInfo.Defs {
		if tn, ok := obj.(*types.TypeName); ok && tn.Name() == g.StructName {
			if _, ok := tn.Type().Underlying().(*types.Struct); ok {
				g.mod = NewModel(g.pkg, tn, g.pkg.CommentsOf(ident), &g.Config)
			}
		}
	}
}

func (g *Generator) Output(cwd string) {
	if g.mod == nil {
		return
	}
	dir, _ := filepath.Rel(cwd, filepath.Dir(g.pkg.GoFiles[0]))
	filename := codegen.GenerateFileSuffix(path.Join(dir, stringsx.LowerSnakeCase(g.StructName)+".go"))
	f := codegen.NewFile(g.pkg.Name, filename)
	g.mod.WriteTo(f)
	_, _ = f.Write()
}

func GetModelByName(g *Generator, name string) *Model {
	if g.StructName == name {
		return g.mod
	}
	return nil
}

type Config struct {
	StructName string
	TableName  string
	Database   string

	WithComments        bool
	WithTableName       bool
	WithTableInterfaces bool
	WithMethods         bool

	FieldPrimaryKey   string
	FieldKeyDeletedAt string
	FieldKeyCreatedAt string
	FieldKeyUpdatedAt string
}

func (g *Config) SetDefault() {
	if g.FieldKeyDeletedAt == "" {
		g.FieldKeyDeletedAt = "DeletedAt"
	}

	if g.FieldKeyCreatedAt == "" {
		g.FieldKeyCreatedAt = "CreatedAt"
	}

	if g.FieldKeyUpdatedAt == "" {
		g.FieldKeyUpdatedAt = "UpdatedAt"
	}

	if g.TableName == "" {
		g.TableName = "t_" + stringsx.LowerSnakeCase(g.StructName)
	}
}
