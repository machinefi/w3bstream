package validator

import (
	"context"
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/rules"
)

func SliceElementRule(ctx context.Context, r *Rule) (Validator, error) {
	raw := []byte("")

	if r.Params != nil {
		if len(r.Params) != 1 {
			return nil, fmt.Errorf(
				"slice should only 1 parameter, but got %d", len(r.Params),
			)
		}
		rule, ok := r.Params[0].(*rules.Rule)
		if !ok {
			return nil, fmt.Errorf("slice parameter should be a valid rule")
		}
		raw = rule.RAW
	}

	compiler := FactoryFromContext(ctx)

	v, err := compiler.Compile(ctx, raw, r.Type.Elem(), nil)
	if err != nil {
		return nil, fmt.Errorf("slice elem %s", err)
	}
	return v, nil
}

func RangeFromUint(min uint64, max *uint64) []*rules.Lit {
	ranges := make([]*rules.Lit, 2)

	if min == 0 && max == nil {
		return nil
	}

	ranges[0] = rules.NewLiteral([]byte(fmt.Sprintf("%d", min)))

	if max != nil {
		if min == *max {
			return []*rules.Lit{ranges[0]}
		}
		ranges[1] = rules.NewLiteral([]byte(fmt.Sprintf("%d", *max)))
	}

	return ranges
}
