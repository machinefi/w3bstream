package codegen

import (
	"fmt"
	"go/token"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func Var(t SnippetType, names ...string) *SnippetField {
	return &SnippetField{Type: t, Names: Idents(names...)}
}

func Type(name string) *NamedType { return &NamedType{Name: Ident(name)} }

func Interface(methods ...IfCanBeIfMethod) *InterfaceType { return &InterfaceType{Methods: methods} }

func Literal(s string) *SnippetLiteral {
	ret := SnippetLiteral(s)
	return &ret
}

func TypeAssert(typ SnippetType, z Snippet) *SnippetTypeAssertExpr {
	return &SnippetTypeAssertExpr{Asserter: z, Type: typ}
}

func Compose(typ SnippetType, elements ...Snippet) *SnippetLiteralCompose {
	return &SnippetLiteralCompose{Type: typ, Elements: elements}
}

func Comments(cmt ...string) *SnippetComments {
	comments := make([]string, 0, len(cmt))
	for _, c := range cmt {
		comments = append(comments, strings.Split(c, "\n")...)
	}
	return &SnippetComments{OneLine: false, Comments: comments}
}

func KeyValue(k, v Snippet) *SnippetKVExpr { return &SnippetKVExpr{K: k, V: v} }

func Ident(s string) *SnippetIdent {
	values := strings.Split(s, ".")

	if !IsValidIdent(values[0]) {
		panic(fmt.Errorf("`%s` is not a valid identifier", values[0]))
	}
	if len(values) == 2 {
		if !IsValidIdent(values[1]) {
			panic(fmt.Errorf("`%s` is not a valid identifier", values[1]))
		}
	}

	ret := SnippetIdent(s)
	return &ret
}

func Idents(s ...string) []*SnippetIdent {
	ret := make([]*SnippetIdent, 0, len(s))
	for _, ident := range s {
		ret = append(ret, Ident(ident))
	}
	return ret
}

func Define(ls ...SnippetCanAddr) *SnippetAssignStmt {
	return &SnippetAssignStmt{Token: token.DEFINE, Ls: ls}
}

func Assign(ls ...SnippetCanAddr) *SnippetAssignStmt {
	return &SnippetAssignStmt{Token: token.ASSIGN, Ls: ls}
}

func AssignWith(tok token.Token, ls ...SnippetCanAddr) *SnippetAssignStmt {
	return &SnippetAssignStmt{Token: tok, Ls: ls}
}

func Inc(v SnippetCanAddr) *SnippetIncExpr { return &SnippetIncExpr{Value: v} }

func Dec(v SnippetCanAddr) *SnippetDecExpr { return &SnippetDecExpr{Value: v} }

func Ref(lead Snippet, refs ...Snippet) *SnippetRefExpr {
	return &SnippetRefExpr{Lead: lead, Refs: refs}
}

func Access(v SnippetCanAddr, index int) *SnippetAccessExpr {
	return &SnippetAccessExpr{V: v, Index: Valuer(index)}
}

func AccessWith(v SnippetCanAddr, index Snippet) *SnippetAccessExpr {
	return &SnippetAccessExpr{V: v, Index: index}
}

func CaseClause(s ...Snippet) *SnippetCaseClause {
	return &SnippetCaseClause{Case: s}
}

func ForRange(ranger Snippet, k, v *SnippetIdent) *SnippetForRangeStmt {
	var kv, vv = AnonymousIdent, AnonymousIdent
	if k != nil {
		kv = *k
	}
	if v != nil {
		vv = *v
	}
	return &SnippetForRangeStmt{Ranger: ranger, K: kv, V: vv}
}

func For(init, cond, post Snippet) *SnippetForStmt {
	return &SnippetForStmt{Init: init, Cond: cond, Post: post}
}

func If(cond Snippet) *SnippetIfStmt { return &SnippetIfStmt{Cond: cond} }

func Switch(cond Snippet) *SnippetSwitchStmt {
	return &SnippetSwitchStmt{Cond: cond}
}

func Select(clauses ...*SnippetCaseClause) *SnippetSelectStmt {
	return &SnippetSelectStmt{Clauses: clauses}
}

func Star(typ SnippetType) *SnippetStarExpr { return &SnippetStarExpr{T: typ} }

func AccessValue(v SnippetCanAddr) *SnippetAccessValueExpr { return &SnippetAccessValueExpr{V: v} }

func Addr(val SnippetCanAddr) *SnippetAddrExpr { return &SnippetAddrExpr{V: val} }

func Paren(s Snippet) *SnippetParenExpr { return &SnippetParenExpr{V: s} }

func Arrow(ch Snippet) *SnippetArrowExpr { return &SnippetArrowExpr{Chan: ch} }

func Casting(ori Snippet, tar Snippet) *SnippetCallExpr { return CallWith(ori, tar) }

func Call(name string, args ...Snippet) *SnippetCallExpr {
	var callee Snippet
	if IsBuiltinFunc(name) {
		callee = SnippetIdent(name)
	} else {
		callee = Ident(name)
	}

	return &SnippetCallExpr{Callee: callee, Args: args}
}

func CallWith(callee Snippet, args ...Snippet) *SnippetCallExpr {
	return &SnippetCallExpr{Callee: callee, Args: args}
}

func CallMakeChan(t SnippetType, length int) *SnippetCallExpr {
	return Call("make", Chan(t), Valuer(length))
}

func Return(s ...Snippet) *SnippetReturnStmt { return &SnippetReturnStmt{Res: s} }

func Func(args ...*SnippetField) *FuncType { return &FuncType{Args: args} }

func Chan(t SnippetType) *ChanType { return &ChanType{T: t, Mode: ChanModeRW} }

func ChanRO(t SnippetType) *ChanType { return &ChanType{T: t, Mode: ChanModeRO} }

func ChanWO(t SnippetType) *ChanType { return &ChanType{T: t, Mode: ChanModeWO} }

func Array(t SnippetType, l int) *ArrayType { return &ArrayType{T: t, Len: l} }

func Slice(t SnippetType) *SliceType { return &SliceType{T: t} }

func Map(k, v SnippetType) *MapType { return &MapType{Tk: k, Tv: v} }

func Ellipsis(t SnippetType) *EllipsisType { return &EllipsisType{T: t} }

func Struct(fields ...*SnippetField) *StructType {
	return &StructType{Fields: fields}
}

func DeclConst(specs ...SnippetSpec) *SnippetTypeDecl {
	return &SnippetTypeDecl{Token: token.CONST, Specs: specs}
}

func DeclVar(specs ...SnippetSpec) *SnippetTypeDecl {
	return &SnippetTypeDecl{Token: token.VAR, Specs: specs}
}

func DeclType(specs ...SnippetSpec) *SnippetTypeDecl {
	return &SnippetTypeDecl{Token: token.TYPE, Specs: specs}
}

func ValueWithAlias(alias FnAlaise) func(interface{}) Snippet {
	return func(v interface{}) Snippet {
		rv := reflect.ValueOf(v)
		rt := reflect.TypeOf(v)

		val := ValueWithAlias(alias)
		typ := TypeWithAlias(alias)

		switch rv.Kind() {
		case reflect.Ptr:
			return Addr(Paren(val(rv.Elem().Interface())))

		case reflect.Struct:
			values := make([]Snippet, 0)
			for i := 0; i < rv.NumField(); i++ {
				fi := rv.Field(i)
				ft := rt.Field(i)

				if !ft.IsExported() {
					continue
				}

				if !IsEmptyValue(fi) {
					values = append(values,
						KeyValue(Ident(ft.Name), val(fi.Interface())))
				}
			}
			return Compose(typ(rt), values...)

		case reflect.Map:
			values := make([]Snippet, 0)
			for _, key := range rv.MapKeys() {
				values = append(values, KeyValue(
					val(key.Interface()),
					val(rv.MapIndex(key).Interface())),
				)
			}
			// to make sure snippet map is ordered
			sort.Slice(values, func(i, j int) bool {
				ik := string(values[i].(*SnippetKVExpr).K.Bytes())
				jk := string(values[j].(*SnippetKVExpr).K.Bytes())
				return ik < jk
			})
			return Compose(typ(rt), values...)

		case reflect.Slice, reflect.Array:
			values := make([]Snippet, 0)
			for i := 0; i < rv.Len(); i++ {
				values = append(values, val(rv.Index(i).Interface()))
			}
			return Compose(typ(rt), values...)

		case reflect.Int, reflect.Uint,
			reflect.Uint8, reflect.Int8,
			reflect.Int16, reflect.Uint16,
			reflect.Int64, reflect.Uint64,
			reflect.Uint32:
			return Literal(fmt.Sprintf("%d", v))

		case reflect.Int32:
			if vr, ok := v.(rune); ok {
				s := strconv.QuoteRune(vr)
				if len(s) == 3 {
					return Literal(s)
				}
			}
			return Literal(fmt.Sprintf("%d", v))

		case reflect.Bool:
			return Literal(strconv.FormatBool(v.(bool)))

		case reflect.Float32:
			return Literal(strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32))

		case reflect.Float64:
			return Literal(strconv.FormatFloat(v.(float64), 'f', -1, 32))

		case reflect.String:
			return Literal(strconv.Quote(v.(string)))

		case reflect.Invalid:
			return SnippetBuiltIn("nil")

		default:
			panic(fmt.Errorf("%v is an unsupported type", v))
		}
	}
}

func TypeWithAlias(aliase FnAlaise) func(reflect.Type) SnippetType {
	return func(t reflect.Type) SnippetType {
		if t.PkgPath() != "" {
			return Type(aliase(t.PkgPath()) + "." + t.Name())
		}

		tof := TypeWithAlias(aliase)
		switch t.Kind() {
		case reflect.Ptr:
			return Star(tof(t.Elem()))
		case reflect.Chan:
			return Chan(tof(t.Elem()))
		case reflect.Array:
			return Array(tof(t.Elem()), t.Len())
		case reflect.Slice:
			return Slice(tof(t.Elem()))
		case reflect.Map:
			return Map(tof(t.Key()), tof(t.Elem()))
		case reflect.Struct:
			fields := make([]*SnippetField, 0)
			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if f.Anonymous {
					fields = append(
						fields,
						Var(tof(f.Type)).WithTag(string(f.Tag)),
					)
				} else {
					fields = append(
						fields,
						Var(tof(f.Type), f.Name).WithTag(string(f.Tag)),
					)
				}
			}
			return Struct(fields...)
		default:
			return BuiltInType(t.String())
		}
	}
}

func ExprWithAlias(alias FnAlaise) func(string, ...interface{}) SnippetExpr {
	val := ValueWithAlias(alias)
	regexpExprHolder := regexp.MustCompile(`(\$\d+)|\?`)

	return func(f string, args ...interface{}) SnippetExpr {
		idx := 0
		return SnippetExpr(
			regexpExprHolder.ReplaceAllStringFunc(
				f,
				func(i string) string {
					arg := args[idx]
					idx++
					if s, ok := arg.(Snippet); ok {
						return Stringify(s)
					}
					return Stringify(val(arg))
				},
			),
		)
	}
}
