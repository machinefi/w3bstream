package transformer

import (
	"context"
	"go/ast"
	"reflect"
	"sync"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	vldterr "github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type Param struct {
	In    string
	Name  string
	Field typesx.StructField
	Type  typesx.Type
	Tags  map[string]reflectx.StructTag
	Loc   []int
}

func (p *Param) FieldValue(rv reflect.Value) reflect.Value {
	rv = reflectx.Indirect(rv)
	n := len(p.Loc)
	fv := rv

	for i := 0; i < n; i++ {
		loc := p.Loc[i]
		fv = fv.Field(loc)

		// last loc should keep ptr value
		if i < n-1 {
			for fv.Kind() == reflect.Ptr {
				// notice the ptr struct ensure only for Ptr Anonymous Field
				if fv.IsNil() {
					fv.Set(reflectx.New(fv.Type()))
				}
				fv = fv.Elem()
			}
		}
	}

	return fv
}

func EachParameter(_ context.Context, t typesx.Type, each func(p *Param) bool) {
	var walk func(tpe typesx.Type, parents ...int)

	walk = func(t typesx.Type, parents ...int) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !ast.IsExported(f.Name()) {
				continue
			}

			loc := append(parents, i)
			flags := reflectx.ParseStructTag(string(f.Tag()))
			display := f.Name()

			tagIn, hasIn := flags["in"]
			tagName, hasName := flags["name"]
			if hasName {
				if name := tagName.Name(); name == "-" {
					continue // skip name:"-"
				} else {
					if name != "" {
						display = name
					}
				}
			}

			if f.Anonymous() && (!hasIn && !hasName) {
				ft := f.Type()
				_, ok := typesx.EncodingTextMarshalerTypeReplacer(ft)
				if !ok {
					for ft.Kind() == reflect.Ptr {
						ft = ft.Elem()
					}
					if ft.Kind() == reflect.Struct {
						walk(ft, loc...)
						continue
					}
				}
			}

			p := &Param{
				In:    tagIn.Name(),
				Name:  display,
				Field: f,
				Type:  f.Type(),
				Tags:  flags,
				Loc:   append([]int{}, loc...),
			}
			if !each(p) {
				break
			}
		}
	}

	walk(t)
}

type FlattenParams struct {
	Params []ReqParam
}

func (FlattenParams) NewValidator(ctx context.Context, typ typesx.Type) (validator.Validator, error) {
	p := &FlattenParams{}
	err := p.CollectParams(ctx, typ)
	return p, err
}

func (FlattenParams) String() string { return "@flatten" }

func (ps *FlattenParams) Validate(v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	errs := vldterr.NewErrorSet()
	rv = reflectx.Indirect(rv)

	for i := range ps.Params {
		pi := ps.Params[i]

		fieldValue := pi.FieldValue(rv)

		if pi.Validator != nil {
			if err := pi.Validator.Validate(fieldValue); err != nil {
				errs.AddErr(err, pi.Name)
			}
		}
	}

	return errs.Err()
}

func (ps *FlattenParams) CollectParams(ctx context.Context, typ typesx.Type) error {
	err := EachReqParam(ctx, typesx.DeRef(typ), func(rp *ReqParam) {
		ps.Params = append(ps.Params, *rp)
	})
	return err
}

type ReqParam struct {
	Param
	Option    Option
	Tsf       Transformer
	Validator validator.Validator
}

func EachReqParam(ctx context.Context, tpe typesx.Type, each func(rp *ReqParam)) error {
	errs := vldterr.NewErrorSet()

	EachParameter(ctx, tpe, func(p *Param) bool {
		rp := &ReqParam{}
		rp.Param = *p

		rp.Option.Name = rp.Name

		if tag, ok := rp.Tags["name"]; ok {
			rp.Option.Omitempty = tag.HasFlag("omitempty")
		}
		if tag, ok := rp.Tags["mime"]; ok {
			rp.Option.MIME = tag.Name()
		}
		if rp.In == "path" {
			rp.Option.Omitempty = false
		}

		switch rp.Type.Kind() {
		case reflect.Array, reflect.Slice:
			elem := rp.Type.Elem()
			if !(elem.PkgPath() == "" && elem.Kind() == reflect.Uint8) {
				rp.Option.Explode = true
			}
		}

		newtsf := func() (Transformer, error) {
			if rp.Option.Explode {
				return NewTransformer(ctx, rp.Type.Elem(), rp.Option)
			}
			return NewTransformer(ctx, rp.Type, rp.Option)
		}

		tsf, err := newtsf()
		if err != nil {
			errs.AddErr(err, rp.Name)
			return false
		}
		rp.Tsf = tsf

		paramVldt, err := NewValidator(ctx, rp.Type, rp.Tags, rp.Option.Omitempty, tsf)
		if err != nil {
			errs.AddErr(err, rp.Name)
			return false
		}
		rp.Validator = paramVldt

		each(rp)

		return true
	})

	if errs.Len() == 0 {
		return nil
	}

	return errs.Err()
}

type ParamAndValue struct {
	Param
	Value reflect.Value
}

type GroupedParams = map[string][]Param

type ckGroupedParamsSet struct{}

// dftGroupedParamsSet stores typesx.Type => GroupedParams
var dftGroupedParamsSet = &sync.Map{}

func GroupedParamSetFromContext(ctx context.Context) *sync.Map {
	if m, ok := ctx.Value(ckGroupedParamsSet{}).(*sync.Map); ok {
		return m
	}
	return dftGroupedParamsSet
}

func WithGroupedParamSet(ctx context.Context, m *sync.Map) context.Context {
	return contextx.WithValue(ctx, ckGroupedParamsSet{}, m)
}

func CollectGroupedParam(ctx context.Context, tpe typesx.Type) GroupedParams {
	if tpe.Kind() != reflect.Struct {
		return nil
	}

	m := GroupedParamSetFromContext(ctx)
	if gp, ok := m.Load(tpe); ok {
		return gp.(GroupedParams)
	}

	gp := GroupedParams{}

	defer func() {
		m.Store(tpe, gp)
	}()

	EachParameter(ctx, tpe, func(p *Param) bool {
		gp[p.In] = append(gp[p.In], *p)
		return true
	})

	return gp
}
