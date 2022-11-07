package validator

import (
	"context"
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/rules"
)

func MapRuleParam(ctx context.Context, r *Rule) (key, elem Validator, err error) {
	if r.Params == nil {
		return
	}

	if len(r.Params) != 2 {
		err = fmt.Errorf(
			"map should only 2 parameter, but got %d", len(r.Params),
		)
		return
	}

	compiler := FactoryFromContext(ctx)

	for i, param := range r.Params {
		switch rule := param.(type) {
		case *rules.Rule:
			switch i {
			case 0:
				key, err = compiler.Compile(ctx, rule.RAW, r.Type.Key())
				if err != nil {
					err = fmt.Errorf("map key %s", err)
					return
				}
			case 1:
				elem, err = compiler.Compile(ctx, rule.RAW, r.Type.Elem())
				if err != nil {
					err = fmt.Errorf("map elem %s", err)
					return
				}
			}
		case *rules.Lit:
			raw := rule.Bytes()
			if len(raw) > 0 {
				err = fmt.Errorf("map parameter should be a valid rule")
				return
			}
			elem, err = compiler.Compile(ctx, raw, r.Type.Elem())
			if err != nil {
				err = fmt.Errorf("map elem %s", err)
				return
			}
		}
	}
	return
}
