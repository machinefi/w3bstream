package middleware

import (
	"context"
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

type AuthPayload struct {
	IdentityID   types.SFID
	IdentityType enums.AccessKeyIdentityType
}

func ParseJwtAuthContentFromContext(ctx context.Context) (*AuthPayload, error) {
	var (
		payload = jwt.AuthFromContext(ctx)
		content []byte
		ret     = &AuthPayload{
			IdentityType: enums.ACCESS_KEY_IDENTITY_TYPE_UNKNOWN,
		}
	)
	switch v := payload.(type) {
	case []byte:
		content = v
	case string:
		content = []byte(v)
	case fmt.Stringer:
		content = []byte(v.String())
	case *models.AccessKey:
		ret.IdentityType = v.IdentityType
		ret.IdentityID = v.IdentityID
		return ret, nil
	default:
		return nil, status.InvalidAuthValue
	}
	err := ret.IdentityID.UnmarshalText(content)
	if err != nil {
		return nil, status.InvalidAuthValue.StatusErr().WithDesc(err.Error())
	}
	return ret, nil
}
