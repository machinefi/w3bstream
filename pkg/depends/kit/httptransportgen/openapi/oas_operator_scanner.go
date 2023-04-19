package openapi

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func NewOperatorScanner(pkg *pkgx.Pkg) *OperatorScanner {
	return &OperatorScanner{
		pkg:              pkg,
		DefScanner:       NewDefScanner(pkg),
		StatusErrScanner: NewStatusErrScanner(pkg),
	}
}

type OperatorScanner struct {
	*DefScanner
	*StatusErrScanner
	pkg *pkgx.Pkg
	ops map[*types.TypeName]*Operator
}

func (os *OperatorScanner) Operator(ctx context.Context, tn *types.TypeName) *Operator {
	if tn == nil {
		return nil
	}

	if tn.Pkg().Path() == PkgPathHttpTspt {
		if tn.Name() == "MetaOperator" || tn.Name() == "GroupOperator" {
			return &Operator{}
		}
	}

	if op, ok := os.ops[tn]; ok {
		return op
	}

	log.FromContext(ctx).Debug(
		"scanning Operator `%s.%s`",
		tn.Pkg().Path(), tn.Name(),
	)

	defer func() {
		if e := recover(); e != nil {
			panic(
				errors.Errorf(
					"scan Operator `%s` failed, panic: %s; calltrace: %s",
					FullTypeName(tn), fmt.Sprint(e), string(debug.Stack()),
				),
			)
		}
	}()

	if t, ok := tn.Type().Underlying().(*types.Struct); ok {
		op := &Operator{
			Tag: TagFromRelativePath(os.pkg.PkgPath, tn),
		}
		os.ScanRouteMeta(op, tn)
		os.ScanRequest(ctx, op, t)
		os.ScanResults(ctx, op, tn)

		// cached scanned
		if os.ops == nil {
			os.ops = map[*types.TypeName]*Operator{}
		}
		os.ops[tn] = op

		return op
	}

	return nil
}

func (os *OperatorScanner) ScanRouterMetaByName(tn *types.TypeName, name string) (string, bool) {
	if tn == nil {
		return "", false
	}

	for _, t := range []types.Type{tn.Type(), types.NewPointer(tn.Type())} {
		_, res, ok := AssertIfByMtdNameAndResCntInPkg(t, os.pkg, name, 1)
		if ok {
			for _, v := range res[0] {
				if v.Value != nil {
					s, err := strconv.Unquote(v.Value.ExactString())
					if err != nil {
						panic(errors.Errorf("%s: %s", err, v.Value))
					}
					return s, true
				}
			}
		}
	}

	return "", false
}

func (os *OperatorScanner) ScanRouteMeta(op *Operator, tn *types.TypeName) {
	op.ID = tn.Name()

	// TODO some router meta from tag, this should be impled in httptransport.RouterMeta
	t := tn.Type().Underlying().(*types.Struct)
	for i := 0; i < t.NumFields(); i++ {
		fi := t.Field(i)
		tags := reflect.StructTag(t.Tag(i))
		if fi.Anonymous() && strings.Contains(fi.Type().String(), PkgPathHttpx+".Method") {
			if path, ok := tags.Lookup("path"); ok {
				vs := strings.Split(path, ",")
				op.Path = vs[0]
				if len(vs) > 0 {
					for i := range vs {
						switch vs[i] {
						case "deprecated":
							op.Deprecated = true
						}
					}
				}
			}
			if basePath, ok := tags.Lookup("basePath"); ok {
				op.BasePath = basePath
			}
			if summary, ok := tags.Lookup("summary"); ok {
				op.Summary = summary
			}
			break
		}
	}

	comments := strings.Split(os.pkg.CommentsOf(os.pkg.IdentOf(tn)), "\n")

	for i := range comments {
		if strings.Contains(comments[i], "@deprecated") {
			op.Deprecated = true
		}
	}

	if op.Summary == "" {
		comments = FilterMarkedLines(comments)
		if comments[0] != "" {
			op.Summary = comments[0]
			if len(comments) > 1 {
				op.Description = strings.Join(comments[1:], "\n")
			}
		}
	}

	if method, ok := os.ScanRouterMetaByName(tn, "Method"); ok {
		op.Method = method
	}

	if path, ok := os.ScanRouterMetaByName(tn, "Path"); ok {
		op.Path = path
	}

	if bath, ok := os.ScanRouterMetaByName(tn, "BasePath"); ok {
		op.BasePath = bath
	}
}

func (os *OperatorScanner) ScanResults(ctx context.Context, op *Operator, tn *types.TypeName) {
	for _, t := range MaybeTypesOfTypeName(tn) {
		mtd, res, ok := AssertIfByMtdNameAndResCntInPkg(t, os.pkg, "Output", 2)
		if ok {
			for _, v := range res[0] {
				vt := v.Type
				if vt == nil {
					continue
				}
				if vt.String() != types.Typ[types.UntypedNil].String() {
					if op.SuccessType != nil && op.SuccessType.String() != v.Type.String() {
						log.FromContext(ctx).Warn(
							errors.Errorf("%s success result must be same struct, "+
								"but got %v, already set %v",
								op.ID, v.Type, op.SuccessType),
						)
					}
					op.SuccessType = vt
					op.SuccessStatus, op.SuccessResponse = os.GetResponse(ctx, vt, v.Expr)
				}
			}
		}
		if os.StatusErrScanner.Type != nil && mtd != nil {
			op.StatusErrors = os.StatusErrScanner.StatusErrorsInFunc(mtd.(*typesx.GoMethod).Func)
			op.StatusErrorSchema = os.DefScanner.GetSchemaByType(ctx, os.StatusErrScanner.Type)
		}
	}
}

func (os *OperatorScanner) ScanRequest(ctx context.Context, op *Operator, t *types.Struct) {
	typesx.EachField(
		typesx.FromGoType(t),
		"name",
		func(f typesx.StructField, name string, omitempty bool) bool {
			loc, _ := reflectx.TagValueAndFlags(f.Tag().Get("in"))

			if loc == "" {
				panic(errors.Errorf("missing tag `in` for %s of %s", f.Name(), op.ID))
			}

			_, flags := reflectx.TagValueAndFlags(f.Tag().Get("name"))

			schema := os.DefScanner.propSchemaByField(
				ctx,
				f.Name(),
				f.Type().(*typesx.GoType).Type,
				f.Tag(),
				flags,
				os.pkg.CommentsOf(os.pkg.IdentOf(f.(*typesx.GoStructField).Var)),
			)

			tsfm, err := transformer.DefaultFactory.NewTransformer(
				context.Background(),
				f.Type(),
				transformer.Option{
					MIME: f.Tag().Get("mime"),
				},
			)
			if err != nil {
				panic(err)
			}

			switch loc {
			case "body":
				body := oas.NewRequestBody("", true)
				body.AddContent(tsfm.Names()[0], oas.NewMediaTypeWithSchema(schema))
				op.SetRequestBody(body)
			case "query":
				op.AddNonBodyParameter(oas.QueryParameter(name, schema, !omitempty))
			case "cookie":
				op.AddNonBodyParameter(oas.CookieParameter(name, schema, !omitempty))
			case "header":
				op.AddNonBodyParameter(oas.HeaderParameter(name, schema, !omitempty))
			case "path":
				op.AddNonBodyParameter(oas.PathParameter(name, schema))
			}

			return true
		},
		"in",
	)
}

func (os *OperatorScanner) FirstConstValueOfFunc(named *types.Named, name string) (interface{}, bool) {
	_, res, ok := AssertIfByMtdNameAndResCntInPkg(types.NewPointer(named), os.pkg, name, 1)
	if ok {
		for _, r := range res[0] {
			if r.IsValue() {
				if v := ConstantValueOf(r.Value); v != nil {
					return v, true
				}
			}
		}
		return nil, true
	}
	return nil, false
}

func (os *OperatorScanner) GetResponse(ctx context.Context, t types.Type, expr ast.Expr) (c int, rsp *oas.Response) {
	rsp = &oas.Response{}

	if t.String() == "error" {
		c = http.StatusNoContent
		return
	}

	ct := ""

	if isHttpxResponse(t) {
		wrapper := func(expr ast.Expr) {
			first := true
			ast.Inspect(expr, func(node ast.Node) bool {
				call, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}
				if first {
					first = false
					v, _ := os.pkg.Eval(call.Args[0])
					t = v.Type
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				switch sel.Sel.Name {
				case "WrapSchema":
					v, _ := os.pkg.Eval(call.Args[0])
					t = v.Type
				case "WrapStatusCode":
					v, _ := os.pkg.Eval(call.Args[0])
					if code, ok := ConstantValueOf(v.Value).(int); ok {
						c = code
					}
					return false
				case "WrapContentType":
					v, _ := os.pkg.Eval(call.Args[0])
					if code, ok := ConstantValueOf(v.Value).(string); ok {
						ct = code
					}
					return false
				}
				return true
			})
		}
		if ident, ok := expr.(*ast.Ident); ok && ident.Obj != nil {
			if stmt, ok := ident.Obj.Decl.(*ast.AssignStmt); ok {
				for _, e := range stmt.Rhs {
					wrapper(e)
				}
			}
		} else {
			wrapper(expr)
		}
	}

	if pointer, ok := t.(*types.Pointer); ok {
		t = pointer.Elem()
	}

	if named, ok := t.(*types.Named); ok {
		if v, ok := os.FirstConstValueOfFunc(named, "ContentType"); ok {
			if s, ok := v.(string); ok {
				ct = s
			}
			if ct == "" {
				ct = "*"
			}
		}
		if v, ok := os.FirstConstValueOfFunc(named, "StatusCode"); ok {
			if i, ok := v.(int64); ok {
				c = int(i)
			}
		}
	}

	if ct == "" {
		ct = httpx.MIME_JSON
	}

	rsp.AddContent(ct, oas.NewMediaTypeWithSchema(os.DefScanner.GetSchemaByType(ctx, t)))

	return
}
