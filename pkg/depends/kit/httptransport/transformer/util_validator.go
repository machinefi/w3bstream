package transformer

import (
	"context"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type MayValidator interface {
	NewValidator(context.Context, typesx.Type) (validator.Validator, error)
}

type WithNamedByTag interface {
	NamedByTag() string
}

func NewValidator(
	ctx context.Context,
	ft typesx.Type,
	tags map[string]reflectx.StructTag,
	omitempty bool,
	tsfm Transformer,
) (validator.Validator, error) {
	if with, ok := tsfm.(WithNamedByTag); ok {
		if tagKey := with.NamedByTag(); tagKey != "" {
			ctx = validator.ContextWithTagKey(ctx, tagKey)
		}
	}

	if t, ok := tsfm.(MayValidator); ok {
		return t.NewValidator(ctx, ft)
	}

	compiler := validator.FactoryFromContext(ctx)
	return compiler.Compile(
		ctx,
		[]byte(tags[validator.TagValidate]),
		ft, func(rule validator.Modifier) {
			if omitempty {
				rule.SetOptional(true)
			}
			if dftv, ok := tags[validator.TagDefault]; ok {
				rule.SetDefaultValue([]byte(dftv))
			}
		},
	)
}
