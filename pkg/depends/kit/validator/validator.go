package validator

import (
	"context"
	"encoding"
	"fmt"
	"reflect"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/rules"
	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/textx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

// Validator can validate input value
type Validator interface {
	Validate(interface{}) error
	String() string
}

type CanValidate interface {
	Validate() error
}

// Creator interface can new a validator
type Creator interface {
	Names() []string
	New(context.Context, *Rule) (Validator, error)
}

// Modifier can change rule's option
type Modifier interface {
	SetOptional(opt bool)
	SetDefaultValue(dft []byte)
	SetErrMsg(err []byte)
}

type Processor = func(rule Modifier)

// Rule with internal error message and type identifier
type Rule struct {
	*rules.Rule
	ErrMsg []byte
	Type   typesx.Type
}

func ParseRuleByType(r []byte, t typesx.Type) (*Rule, error) {
	rule := &rules.Rule{}
	if len(r) != 0 {
		parsed, err := rules.ParseRaw(r)
		if err != nil {
			return nil, err
		}
		rule = parsed
	}
	return &Rule{Rule: rule, Type: t}, nil
}

func ParseRuleStringByType(r string, t typesx.Type) (*Rule, error) {
	return ParseRuleByType([]byte(r), t)
}

func MustParseRuleByType(r []byte, t typesx.Type) *Rule {
	rule, err := ParseRuleByType(r, t)
	if err != nil {
		panic(err)
	}
	return rule
}

func MustParseRuleStringByType(r string, t typesx.Type) *Rule {
	return MustParseRuleByType([]byte(r), t)
}

func (r *Rule) SetOptional(opt bool) { r.Optional = opt }

func (r *Rule) SetDefaultValue(dft []byte) { r.DftValue = dft }

func (r *Rule) SetErrMsg(msg []byte) { r.ErrMsg = msg }

func (r *Rule) String() string {
	return typesx.FullTypeName(r.Type) + string(r.Rule.Bytes())
}

type Factory interface {
	Compile(context.Context, []byte, typesx.Type, ...Processor) (Validator, error)
}

type ckCompiler struct{}

func ContextWithFactory(ctx context.Context, c Factory) context.Context {
	return contextx.WithValue(ctx, ckCompiler{}, c)
}

func FactoryFromContext(ctx context.Context) Factory {
	if f, ok := ctx.Value(ckCompiler{}).(Factory); ok {
		return f
	}
	return DefaultFactory
}

type factory struct {
	set map[string]Creator
}

var DefaultFactory = NewFactory()

func NewFactory() *factory { return &factory{set: make(map[string]Creator)} }

func (f *factory) Register(creators ...Creator) {
	for i := range creators {
		for _, name := range creators[i].Names() {
			f.set[name] = creators[i]
		}
	}
}

func (f *factory) MustCompile(ctx context.Context, rule []byte, t typesx.Type, processors ...Processor) Validator {
	v, err := f.Compile(ctx, rule, t, processors...)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *factory) Compile(ctx context.Context, rule []byte, t typesx.Type, processors ...Processor) (Validator, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(rule) == 0 {
		if _, ok := typesx.EncodingTextMarshalerTypeReplacer(t); !ok {
			switch typesx.DeRef(t).Kind() {
			case reflect.Struct:
				rule = []byte("@struct")
			case reflect.Slice:
				rule = []byte("@slice")
			case reflect.Map:
				rule = []byte("@map")
			}
		}
	}
	r, err := ParseRuleByType(rule, t)
	if err != nil {
		return nil, err
	}
	for _, proc := range processors {
		if proc != nil {
			proc(r)
		}
	}
	creator, ok := f.set[r.Name]
	if !ok && len(rule) > 0 {
		return nil, fmt.Errorf("%s not match any validator", r.Name)
	}
	return NewLoader(creator).New(ContextWithFactory(ctx, f), r)
}

type PreprocessStage int

const (
	PreprocessSkip PreprocessStage = iota
	PreprocessString
	PreprocessPtr
)

// Loader load from creator
type Loader struct {
	Creator Creator
	Validator
	PreprocessStage

	DftValue []byte
	Optional bool
	ErrMsg   []byte
}

func NewLoader(c Creator) *Loader { return &Loader{Creator: c} }

func normalize(typ typesx.Type) (typesx.Type, PreprocessStage) {
	if t, ok := typesx.EncodingTextMarshalerTypeReplacer(typ); ok {
		return t, PreprocessString
	}
	if typ.Kind() == reflect.Ptr {
		return typesx.DeRef(typ), PreprocessPtr
	}
	return typ, PreprocessSkip
}

func (l *Loader) String() string {
	if l.Validator != nil {
		v := l.Validator.String()

		if l.Optional {
			if l.DftValue != nil {
				return v + " = " + string(rules.SingleQuote(l.DftValue))
			}
			return v + "?"
		}

		return v
	}
	return "nil"
}

func (l *Loader) New(ctx context.Context, rule *Rule) (Validator, error) {
	ret := NewLoader(l.Creator)

	ret.Optional = rule.Optional
	ret.DftValue = rule.DftValue
	ret.ErrMsg = rule.ErrMsg

	typ := rule.Type
	rule.Type, ret.PreprocessStage = normalize(rule.Type)

	if l.Creator != nil {
		v, err := l.Creator.New(ctx, rule)
		if err != nil {
			return nil, err
		}
		ret.Validator = v

		if ret.DftValue != nil {
			if rv, ok := typesx.TryNew(typ); ok {
				if err := textx.UnmarshalText(rv, ret.DftValue); err != nil {
					return nil, fmt.Errorf(
						"default value `%s` can not unmarshal to %s: %s",
						ret.DftValue, typ, err,
					)
				}
				if err := ret.Validate(rv); err != nil {
					return nil, fmt.Errorf(
						"default value `%s` is not a valid value of %s: %s",
						ret.DftValue, v, err,
					)
				}
			}
		}
	}

	return ret, nil
}

func (l *Loader) Validate(v interface{}) error {
	err := l.validate(v)
	if err == nil {
		return nil
	}
	if l.ErrMsg != nil && len(l.ErrMsg) != 0 {
		return fmt.Errorf(string(l.ErrMsg))
	}
	return err
}

func (l *Loader) validate(v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if reflectx.IsEmptyValue(rv) {
		if !l.Optional {
			return errors.MissingRequiredFieldError{}
		}

		if l.DftValue != nil && rv.CanSet() {
			if err := textx.UnmarshalText(rv, l.DftValue); err != nil {
				return fmt.Errorf("unmarshal default value failed")
			}
		}
		// empty value should not to validate
		return nil
	}

	if l.Validator == nil {
		return nil
	}

	if l.PreprocessStage == PreprocessString {
		// make sure value over reflect.Value
		if rv.CanInterface() {
			v = rv.Interface()
		}

		if marshaller, ok := v.(encoding.TextMarshaler); ok {
			data, err := marshaller.MarshalText()
			if err != nil {
				return err
			}
			return l.Validator.Validate(string(data))
		}
	}

	if rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	return l.Validator.Validate(reflectx.Indirect(rv))
}
