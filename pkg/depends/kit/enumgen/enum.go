package enumgen

import (
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

type Enum struct {
	Path    string
	Name    string
	Options Options
}

func NewEnum(name string, options Options) *Enum {
	enum := &Enum{Options: options}
	parts := strings.Split(name, ".")
	switch len(parts) {
	case 1:
		enum.Path = ""
		enum.Name = parts[0]
	default:
		enum.Path = strings.Join(parts[0:len(parts)-1], ".")
		enum.Name = parts[len(parts)-1]
	}
	return enum
}

func (e Enum) ConstUnknown() Snippet {
	return Ident(stringsx.UpperSnakeCase(e.Name) + "_UNKNOWN")
}

func (e Enum) ConstName(value string) Snippet {
	return Ident(stringsx.UpperSnakeCase(e.Name) + "__" + value)
}

func (e Enum) VarInvalidError() Snippet { return Ident("Invalid" + e.Name) }

func (e Enum) Receiver(name string) *SnippetField { return Var(Type(e.Name), name) }

func (e Enum) StarReceiver(name string) *SnippetField { return Var(Star(Type(e.Name)), name) }

func (e Enum) WriteToFile(f *File) {
	f.WriteSnippet(
		e.Errors(f),
		e.StringParser(f),
		e.LabelParser(f),
		e.Integer(f),
		e.Stringer(f),
		e.Labeler(f),
		e.TypeName(f),
		e.ConstValues(f),
		e.TextMarshaler(f),
		e.TextUnmarshaler(f),
		e.Scanner(f),
		e.Valuer(f),
	)
}

func (e Enum) Errors(f *File) Snippet {
	return Exprer(
		`var ? = ?("invalid ? type")`,
		e.VarInvalidError(),
		Ident(f.Use("errors", "New")),
		Ident(e.Name),
	)
}

func (e Enum) StringParser(f *File) Snippet {
	clauses := []*SnippetCaseClause{
		CaseClause().Do(Return(e.ConstUnknown(), e.VarInvalidError())),
		CaseClause(f.Value("")).Do(Return(e.ConstUnknown(), Nil)),
	}

	for _, o := range e.Options {
		clauses = append(
			clauses,
			CaseClause(f.Value(*o.Str)).Do(Return(e.ConstName(*o.Str), Nil)),
		)
	}

	return Func(Var(String, "s")).
		Named("Parse"+e.Name+"FromString").
		Return(Var(Type(e.Name)), Var(Error)).
		Do(Switch(Ident("s")).When(clauses...))
}

func (e Enum) LabelParser(f *File) Snippet {
	clauses := []*SnippetCaseClause{
		CaseClause().Do(Return(e.ConstUnknown(), e.VarInvalidError())),
		CaseClause(f.Value("")).Do(Return(e.ConstUnknown(), Nil)),
	}

	for _, o := range e.Options {
		clauses = append(
			clauses,
			CaseClause(f.Value(o.Label)).Do(Return(e.ConstName(*o.Str), Nil)),
		)
	}

	return Func(Var(String, "s")).
		Named("Parse"+e.Name+"FromLabel").
		Return(Var(Type(e.Name)), Var(Error)).
		Do(Switch(Ident("s")).When(clauses...))
}

func (e Enum) Stringer(f *File) Snippet {
	clauses := []*SnippetCaseClause{
		CaseClause().Do(Return(f.Value("UNKNOWN"))),
		CaseClause(e.ConstUnknown()).Do(Return(f.Value(""))),
	}

	for _, o := range e.Options {
		clauses = append(
			clauses,
			CaseClause(e.ConstName(*o.Str)).Do(Return(f.Value(*o.Str))),
		)
	}

	return Func().
		MethodOf(e.Receiver("v")).
		Named("String").
		Return(Var(String)).
		Do(Switch(Ident("v")).When(clauses...))
}

func (e Enum) Integer(f *File) Snippet {
	return Func().
		MethodOf(e.Receiver("v")).
		Named("Int").
		Return(Var(Int)).
		Do(Return(Casting(Int, Ident("v"))))
}

func (e Enum) Labeler(f *File) Snippet {
	clauses := []*SnippetCaseClause{
		CaseClause().Do(Return(f.Value("UNKNOWN"))),
		CaseClause(e.ConstUnknown()).Do(Return(f.Value(""))),
	}

	for _, o := range e.Options {
		clauses = append(
			clauses,
			CaseClause(e.ConstName(*o.Str)).Do(Return(f.Value(o.Label))),
		)
	}

	return Func().
		MethodOf(e.Receiver("v")).
		Named("Label").
		Return(Var(String)).
		Do(Switch(Ident("v")).When(clauses...))
}

func (e Enum) TypeName(f *File) Snippet {
	return Func().
		MethodOf(e.Receiver("v")).
		Named("TypeName").
		Return(Var(String)).
		Do(Return(f.Value(e.Path + "." + e.Name)))
}

func (e Enum) ConstValues(f *File) Snippet {
	typ := Slice(Type(f.Use(PkgPath, IntStringerName)))
	return Func().
		MethodOf(e.Receiver("v")).
		Named("ConstValues").
		Return(Var(typ)).
		Do(Return(func() Snippet {
			lst := []interface{}{typ}
			holder := "?"
			for i, o := range e.Options {
				if i > 0 {
					holder += ", ?"
				}
				lst = append(lst, e.ConstName(*o.Str))
			}
			return f.Expr("?{"+holder+"}", lst...)
		}()))
}

func (e Enum) TextMarshaler(f *File) Snippet {
	return Func().
		MethodOf(e.Receiver("v")).
		Named("MarshalText").
		Return(Var(Slice(Byte)), Var(Error)).
		Do(
			Define(Ident("s")).By(Ref(Ident("v"), Call("String"))),
			If(f.Expr(`s == "UNKNOWN"`)).
				Do(Return(Nil, e.VarInvalidError())),
			Return(Casting(Slice(Byte), Ident("s")), Nil),
		)

}

func (e Enum) TextUnmarshaler(f *File) Snippet {
	return Func(Var(Slice(Byte), "data")).
		MethodOf(e.StarReceiver("v")).
		Named("UnmarshalText").
		Return(Var(Error)).
		Do(
			Define(Ident("s")).
				By(
					Casting(String, Call(f.Use("bytes", "ToUpper"), Ident("data"))),
				),
			Define(Ident("val"), Ident("err")).
				By(
					CallWith(
						f.Expr("Parse?FromString", Ident(e.Name)),
						Ident("s"),
					),
				),
			If(Exprer("err != nil")).Do(Return(Ident("err"))),
			Assign(AccessValue(Ident("v"))).By(Ident("val")),
			Return(Nil),
		)
}

func (e Enum) Scanner(f *File) Snippet {
	return Func(Var(Interface(), "src")).
		MethodOf(e.StarReceiver("v")).
		Named("Scan").
		Return(Var(Error)).
		Do(
			Define(Ident("offset")).By(Valuer(0)),
			Define(Ident("o"), Ident("ok")).
				By(
					Ref(
						Casting(Interface(), Ident("v")),
						Paren(Exprer(f.Use(PkgPath, ValueOffsetName))),
					),
				),
			If(Ident("ok")).
				Do(
					Assign(Ident("offset")).
						By(Ref(Ident("o"), Call("Offset"))),
				),
			Define(Ident("i"), Ident("err")).
				By(
					Call(
						f.Use(PkgPath, ScanIntEnumStringerName),
						Ident("src"), Ident("offset"),
					),
				),
			If(Exprer("err != nil")).
				Do(Return(Ident("err"))),
			Assign(AccessValue(Ident("v"))).
				By(
					Casting(Type(e.Name), Ident("i")),
				),
			Return(Nil),
		)
}

func (e Enum) Valuer(f *File) Snippet {
	return Func().
		MethodOf(e.Receiver("v")).
		Named("Value").
		Return(
			Var(Type(f.Use("database/sql/driver", "Value"))),
			Var(Error)).
		Do(
			Define(Ident("offset")).By(Valuer(0)),
			Define(Ident("o"), Ident("ok")).
				By(
					Ref(
						Casting(Interface(), Ident("v")),
						Paren(Exprer(f.Use(PkgPath, ValueOffsetName))),
					),
				),
			If(Ident("ok")).
				Do(
					Assign(Ident("offset")).
						By(Ref(Ident("o"), Call("Offset"))),
				),
			Return(Exprer(`int64(v) + int64(offset)`), Nil),
		)
}

var (
	PkgPath                 string
	IntStringerName         = "IntStringerEnum"
	ValueOffsetName         = "ValueOffset"
	ScanIntEnumStringerName = "ScanIntEnumStringer"
)

func init() {
	_, current, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(current), "../enum")
	PkgPath = must.String(pkgx.PkgIdByPath(dir))
}
