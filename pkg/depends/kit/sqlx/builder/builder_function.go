package builder

import "context"

type Function struct {
	name  string
	exprs []SqlExpr
}

func Func(name string, es ...SqlExpr) *Function {
	if name == "" {
		return nil
	}
	return &Function{name: name, exprs: es}
}

func (f *Function) IsNil() bool { return f == nil || f.name == "" }

func (f *Function) Ex(ctx context.Context) *Ex {
	e := Expr(f.name)
	e.WriteGroup(func(e *Ex) {
		if len(f.exprs) == 0 {
			e.WriteQueryByte('*')
		}
		for i := range f.exprs {
			if i > 0 {
				e.WriteQueryByte(',')
			}
			e.WriteExpr(f.exprs[i])
		}
	})
	return e.Ex(ctx)
}

func Count(es ...SqlExpr) *Function {
	if len(es) == 0 {
		return Func("COUNT", Expr("1"))
	}
	return Func("COUNT", es...)
}

func Avg(es ...SqlExpr) *Function { return Func("AVG", es...) }

func Distinct(es ...SqlExpr) *Function { return Func("DISTINCT", es...) }

func Min(es ...SqlExpr) *Function { return Func("MIN", es...) }

func Max(es ...SqlExpr) *Function { return Func("MAX", es...) }

func First(es ...SqlExpr) *Function { return Func("FIRST", es...) }

func Last(es ...SqlExpr) *Function { return Func("LAST", es...) }

func Sum(es ...SqlExpr) *Function { return Func("SUM", es...) }
