package jwt

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
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

	if WithAnonymousPublisherFn != nil {
		tok, err = WithAnonymousPublisherFn(ctx, tok)
	}

	ok = false
	if BuiltInTokenValidateFn != nil {
		pl, err, ok = BuiltInTokenValidateFn(ctx, tok)
	}
	if !ok {
		var claims *Claims
		if claims, err = jwt.ParseToken(tok); err == nil {
			pl = claims.Payload
		}
	}

	if err != nil {
		return nil, err
	}

	if WithPermissionFn != nil && !WithPermissionFn(pl) {
		return nil, ErrNoPermission
	}

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

var BuiltInTokenValidateFn func(context.Context, string) (interface{}, error, bool)

func SetBuiltInTokenFn(f func(context.Context, string) (interface{}, error, bool)) {
	BuiltInTokenValidateFn = f
}

var WithPermissionFn func(interface{}) bool

func SetWithPermissionFn(f func(interface{}) bool) {
	WithPermissionFn = f
}

var WithAnonymousPublisherFn func(context.Context, string) (string, error)

func SetWithAnonymousPublisherFn(f func(context.Context, string) (string, error)) {
	WithAnonymousPublisherFn = f
}
