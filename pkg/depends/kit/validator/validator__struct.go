package validator

import (
	"context"
	"go/ast"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func NewStructValidator(tagKey string) *Struct {
	return &Struct{
		tag:        tagKey,
		validators: map[string]Validator{},
	}
}

type ckStructTagKey struct{}

func ContextWithTagKey(ctx context.Context, namedTagKey string) context.Context {
	return context.WithValue(ctx, ckStructTagKey{}, namedTagKey)
}

func TagKeyFromContext(ctx context.Context) string {
	v := ctx.Value(ckStructTagKey{})
	if v != nil {
		if namedTagKey, ok := v.(string); ok {
			return namedTagKey
		}
	}
	return ""
}

type Struct struct {
	tag        string
	validators map[string]Validator // validators: field name => validator
}

func init() { DefaultFactory.Register(&Struct{}) }

func (Struct) Names() []string { return []string{"struct"} }

func (vs *Struct) Validate(v interface{}) error {
	switch rv := v.(type) {
	case reflect.Value:
		return vs.ValidateReflectValue(rv)
	default:
		return vs.ValidateReflectValue(reflect.ValueOf(v))
	}
}

func (vs *Struct) ValidateReflectValue(rv reflect.Value) error {
	errs := errors.NewErrorSet()
	vs.validate(rv, errs)
	return errs.Err()
}

func (vs *Struct) validate(rv reflect.Value, errs *errors.ErrorSet) {
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Type().Field(i)
		fv := rv.Field(i)
		name, _, exists := typesx.FieldDisplayName(f.Tag, vs.tag, f.Name)

		if !ast.IsExported(f.Name) || name == "-" {
			continue
		}

		ft_ := reflectx.DeRef(f.Type)
		isSub := ft_.Kind() == reflect.Struct

		if f.Anonymous && isSub && !exists {
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				fv = reflectx.New(ft_)
			}
			vs.validate(fv, errs)
			continue
		}

		if validator, ok := vs.validators[f.Name]; ok {
			err := validator.Validate(fv)
			errs.AddErr(err, name)
		}
	}
}

func (vs *Struct) TypeCheck(r *Rule) error {
	if r.Type.Kind() != reflect.Struct {
		return errors.NewUnsupportedTypeError(r.String(), vs.String())
	}
	return nil
}

const (
	TagValidate = "validate"
	TagDefault  = "default"
	TagErrMsg   = "errMsg"
)

func (vs *Struct) New(ctx context.Context, r *Rule) (Validator, error) {
	if err := vs.TypeCheck(r); err != nil {
		return nil, err
	}

	tag := TagKeyFromContext(ctx)
	if r.Rule != nil && len(r.Params) > 0 {
		tag = string(r.Params[0].Bytes())
	}
	if tag == "" {
		tag = vs.tag
	}

	var (
		ret      = NewStructValidator(tag)
		errs     = errors.NewErrorSet()
		compiler = FactoryFromContext(ctx)
	)
	ctx = ContextWithTagKey(ctx, ret.tag)

	typesx.EachField(
		r.Type,
		ret.tag,
		func(field typesx.StructField, display string, omitempty bool) bool {
			tagv := field.Tag().Get(TagValidate)

			if tagv == "" && typesx.DeRef(field.Type()).Kind() == reflect.Struct {
				_, ok := typesx.EncodingTextMarshalerTypeReplacer(field.Type())
				if !ok {
					tagv = ret.String()
				}
			}

			sub, err := compiler.Compile(
				ContextWithTagKey(ctx, tag),
				[]byte(tagv),
				field.Type(),
				func(rule Modifier) {
					if omitempty {
						rule.SetOptional(omitempty)
					}
					if dftv, ok := field.Tag().Lookup(TagDefault); ok {
						rule.SetDefaultValue([]byte(dftv))
					}
					if errMsg, ok := field.Tag().Lookup(TagErrMsg); ok {
						rule.SetErrMsg([]byte(errMsg))
					}
				},
			)

			if err != nil {
				errs.AddErr(err, field.Name())
				return true
			}

			if sub != nil {
				ret.validators[field.Name()] = sub
			}
			return true
		},
	)

	return ret, errs.Err()
}

func (vs *Struct) String() string {
	return "@" + vs.Names()[0] + "<" + vs.tag + ">"
}
