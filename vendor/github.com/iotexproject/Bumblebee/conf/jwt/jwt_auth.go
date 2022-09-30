package jwt

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/iotexproject/Bumblebee/x/contextx"
	"github.com/iotexproject/Bumblebee/x/misc/must"
)

type Auth struct {
	AuthInQuery  string `name:"authorization,omitempty" in:"query"  validate:"@string[1,]"`
	AuthInHeader string `name:"Authorization,omitempty" in:"header" validate:"@string[1,]"`
}

func (r Auth) ContextKey() interface{} { return keyAuth{} }

func (r Auth) Output(ctx context.Context) (pl interface{}, err error) {
	jwt, ok := ConfFromContext(ctx)
	if !ok {
		return nil, ErrEmptyJwtContext
	}

	av := r.AuthInQuery
	if av == "" {
		av = r.AuthInHeader
	}
	tok := strings.TrimSpace(strings.Replace(av, "Bearer", " ", 1))
	if BuiltInTokenFn != nil {
		if pl, ok = BuiltInTokenFn(tok); ok {
			return pl, nil
		}
	}

	claims, err := jwt.ParseToken(tok)
	if err != nil {
		return nil, err
	}

	if WithPermissionFn != nil && !WithPermissionFn(claims) {
		return nil, ErrNoPermission
	}

	pl = claims.Payload
	return
}

type keyConf struct{}

func WithConfContext(jwt *Jwt) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, keyConf{}, jwt)
	}
}

func ConfFromContext(ctx context.Context) (*Jwt, bool) {
	j, ok := ctx.Value(keyConf{}).(*Jwt)
	return j, ok
}

func MustConfFromContext(ctx context.Context) *Jwt {
	j, ok := ctx.Value(keyConf{}).(*Jwt)
	must.BeTrue(ok)
	return j
}

type keyAuth struct{}

func AuthFromContext(ctx context.Context) interface{} {
	return ctx.Value(keyAuth{})
}

var (
	ErrEmptyJwtContext = errors.New("empty jwt context")
	ErrNoPermission    = errors.New("no permission")
)

var BuiltInTokenFn func(string) (interface{}, bool)

func SetBuiltInTokenFn(f func(string) (interface{}, bool)) {
	BuiltInTokenFn = f
}

var WithPermissionFn func(*Claims) bool

func SetWithPermissionFn(f func(*Claims) bool) {
	WithPermissionFn = f
}
