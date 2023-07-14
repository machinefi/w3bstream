package middleware_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis/middleware"
	"github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/models"
)

func TestParseJwtAuthContentFromContext(t *testing.T) {
	var (
		ctxbg = context.Background()
		key   = (&jwt.Auth{}).ContextKey()
		cases = []*struct {
			name        string
			inputCtxVal interface{}
			outputVal   *middleware.AuthPayload
			outputErr   error
		}{
			{"#WrongContent", "wrong_content", nil, status.InvalidAuthValue},
			{"#WrongContentType", 1, nil, status.InvalidAuthValue},
			{"#String", "1", &middleware.AuthPayload{1, 0}, nil},
			{"#Byte", []byte("1"), &middleware.AuthPayload{1, 0}, nil},
			{"#Stringer", time.Now(), nil, status.InvalidAuthValue},
			{"#AccessKey", &models.AccessKey{}, &middleware.AuthPayload{0, 0}, nil},
		}
	)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			output, err := middleware.ParseJwtAuthContentFromContext(context.WithValue(ctxbg, key, c.inputCtxVal))
			if err == nil {
				NewWithT(t).Expect(c.outputErr).To(BeNil())
				NewWithT(t).Expect(*output).To(Equal(*c.outputVal))
			}
		})
	}
}
