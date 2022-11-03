package enumgen

import (
	"go/types"
	"log"
	"path"
	"path/filepath"
	"sort"

	"golang.org/x/tools/go/packages"

	"github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

type Generator struct {
	pkg *pkgx.Pkg
	*Scanner
	enums map[*types.TypeName]*Enum
}

func New(pkg *pkgx.Pkg) *Generator {
	return &Generator{
		pkg:     pkg,
		Scanner: NewScanner(pkg),
		enums:   map[*types.TypeName]*Enum{},
	}
}

func (g *Generator) Scan(names ...string) {
	for _, name := range names {
		tn := g.pkg.TypeName(name)
		if tn == nil {
			continue
		}
		opts, ok := g.Scanner.Options(tn)
		if ok && opts[0].Int != nil && opts[0].Str != nil {
			sort.Slice(opts, func(i, j int) bool {
				return *opts[i].Int < *opts[j].Int
			})
			g.enums[tn] = NewEnum(tn.Pkg().Path()+"."+tn.Name(), opts)
		}
	}
}

func (g Generator) Output(cwd string) {
	for tn, enum := range g.enums {
		dir, _ := filepath.Rel(
			cwd,
			must.String(pkgx.PkgPathByPath(tn.Pkg().Path(), packages.NeedName, packages.NeedFiles)),
		)
		filename := codegen.GenerateFileSuffix(
			path.Join(dir, stringsx.LowerSnakeCase(enum.Name)+".go"))
		f := codegen.NewFile(tn.Pkg().Name(), filename)
		enum.WriteToFile(f)

		if _, err := f.Write(); err != nil {
			log.Printf("%s generate failed: %v", filename, err)
		}
	}
}

func GetEnumByName(g *Generator, name string) *Enum {
	for tn, enum := range g.enums {
		if tn.Name() == name {
			return enum
		}
	}
	return nil
}
