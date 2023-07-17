package jwt_test

import (
	"context"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	. "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
)

func TestJwt(t *testing.T) {
	conf := &Jwt{
		Issuer:  "jwt_unit_test",
		ExpIn:   types.Duration(time.Second),
		SignKey: "any",
	}

	t.Run("#JwtConf", func(t *testing.T) {
		t.Run("#Init", func(t *testing.T) {
			c := *conf
			c.SignKey = ""
			c.Init()
			NewWithT(t).Expect(c.SignKey).To(Equal("xxxx"))

			c.ExpIn = 0
			c.Init()
			NewWithT(t).Expect(c.ExpIn).To(Equal(types.Duration(time.Hour)))
		})

		t.Run("#ExpiresAt", func(t *testing.T) {
			c := *conf
			NewWithT(t).Expect(c.ExpiresAt()).NotTo(BeNil())

			c.ExpIn = 0
			NewWithT(t).Expect(c.ExpiresAt()).To(BeNil())
		})

		t.Run("#GeneratingAndParsing", func(t *testing.T) {
			conf.ExpIn = types.Duration(time.Second * 2)
			payload := interface{}("any")
			tok, err := conf.GenerateTokenByPayload(payload)
			NewWithT(t).Expect(err).To(BeNil())
			NewWithT(t).Expect(tok).NotTo(BeEmpty())

			t.Run("#Success", func(t *testing.T) {
				claim, err := conf.ParseToken(tok)
				NewWithT(t).Expect(err).To(BeNil())

				NewWithT(t).Expect(claim.Payload).To(Equal(payload))
			})

			t.Run("#Failed", func(t *testing.T) {
				t.Run("#TokenExpired", func(t *testing.T) {
					time.Sleep(2 * time.Second)

					_, err = conf.ParseToken(tok)
					NewWithT(t).Expect(err).NotTo(BeNil())

					ve, ok := err.(*statusx.StatusErr)
					NewWithT(t).Expect(ok).To(BeTrue())
					NewWithT(t).Expect(ve.Key).To(Equal(InvalidToken.Key()))
				})

				t.Run("#ParseWithClaimFailed", func(t *testing.T) {
					_, err = conf.ParseToken("not equal token gen before")
					NewWithT(t).Expect(err).NotTo(BeNil())
				})
			})
		})
	})

	auth := &Auth{
		AuthInQuery:  "",
		AuthInHeader: "",
	}
	ctx := WithConfContext(conf)(context.Background())

	fromCtx := MustConfFromContext(ctx)
	NewWithT(t).Expect(fromCtx).To(Equal(conf))

	fromCtx, _ = ConfFromContext(ctx)
	NewWithT(t).Expect(fromCtx).To(Equal(conf))

	t.Run("#AuthMiddleware", func(t *testing.T) {
		payload := interface{}("any")
		t.Run("#Output", func(t *testing.T) {
			tok, _ := conf.GenerateTokenWithoutExpByPayload(payload)
			t.Run("#InQuery", func(t *testing.T) {
				auth.AuthInHeader = ""
				auth.AuthInQuery = tok
				pl, err := auth.Output(ctx)
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(pl).To(Equal(payload))
			})
			t.Run("#InHeader", func(t *testing.T) {
				auth.AuthInHeader = tok
				auth.AuthInQuery = ""
				pl, err := auth.Output(ctx)
				NewWithT(t).Expect(err).To(BeNil())
				NewWithT(t).Expect(pl).To(Equal(payload))
			})

			t.Run("#WithBuiltinValidator", func(t *testing.T) {
				anyError := errors.New("any")
				cases := []*struct {
					name          string
					builtin       func(_ context.Context, _ string) (interface{}, error, bool)
					expectPayload interface{}
					expectError   error
				}{
					{
						name:          "$WithoutBuiltinFn",
						builtin:       nil,
						expectPayload: payload,
					},
					{
						name: "$CannotBeValidated",
						builtin: func(_ context.Context, _ string) (interface{}, error, bool) {
							return nil, nil, false
						},
						expectPayload: payload,
					},
					{
						name: "$ValidatorReturnsError",
						builtin: func(_ context.Context, _ string) (interface{}, error, bool) {
							return nil, anyError, true
						},
						expectError: anyError,
					},
					{
						name: "$ValidatorSuccessValidating",
						builtin: func(_ context.Context, _ string) (interface{}, error, bool) {
							return payload.(string), nil, true
						},
						expectPayload: payload,
					},
				}
				for _, c := range cases {
					t.Run(c.name, func(t *testing.T) {
						SetBuiltInTokenFn(c.builtin)

						pl, err := auth.Output(ctx)

						if c.expectError != nil {
							NewWithT(t).Expect(err).To(Equal(c.expectError))
						} else {
							NewWithT(t).Expect(err).To(BeNil())
							NewWithT(t).Expect(pl).To(Equal(c.expectPayload))
						}
					})
				}
				SetBuiltInTokenFn(nil)
			})

			t.Run("#ParseTokenFailed", func(t *testing.T) {
				SetBuiltInTokenFn(nil)
				if runtime.GOOS == `darwin` {
					return
				}
				patch := gomonkey.ApplyMethod(
					reflect.TypeOf(&Jwt{}),
					"ParseToken",
					func(*Jwt, string) (*Claims, error) {
						return nil, errors.New("any")
					},
				)
				defer patch.Reset()
				_, err := auth.Output(ctx)
				NewWithT(t).Expect(err).NotTo(BeNil())
			})

			t.Run("#WithPermissionFn", func(t *testing.T) {
				cases := []*struct {
					name          string
					fn            func(interface{}) bool
					expectError   error
					expectPayload interface{}
				}{
					{
						name:          "$WithoutPermissionFn",
						fn:            nil,
						expectPayload: payload,
					},
					{
						name:          "$PermissionOK",
						fn:            func(interface{}) bool { return true },
						expectPayload: payload,
					},
					{
						name:        "$PermissionDenied",
						fn:          func(interface{}) bool { return false },
						expectError: ErrNoPermission,
					},
				}

				for _, c := range cases {
					t.Run(c.name, func(t *testing.T) {
						SetWithPermissionFn(c.fn)

						pl, err := auth.Output(ctx)

						if c.expectError != nil {
							NewWithT(t).Expect(err).To(Equal(c.expectError))
						} else {
							NewWithT(t).Expect(err).To(BeNil())
							NewWithT(t).Expect(pl).To(Equal(c.expectPayload))
						}
					})
				}
				SetWithPermissionFn(nil)
			})
		})
	})
}
