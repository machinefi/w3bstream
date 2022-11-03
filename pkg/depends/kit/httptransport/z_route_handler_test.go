package httptransport_test

import (
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/mock"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/testutil/httptransporttestutil/server/cmd/app/routes"
)

var (
	factory = NewRequestTsfmFactory(nil, nil)
	meta    = &ServiceMeta{Name: "service-test", Version: "1.0.0"}
)

func init() {
	factory.SetDefault()
}

func TestHttpRouteHandler(t *testing.T) {
	t.Run("Redirect", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(routes.Redirect{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		operator := routes.Redirect{}
		req, err := factory.NewRequest(operator.Method(), "/", operator)
		NewWithT(t).Expect(err).To(BeNil())

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).To(Equal(`HTTP/0.0 302 Found
Content-Type: text/html; charset=utf-8
Location: /other
X-Meta: service-test@1.0.0/Redirect

<a href="/other">Found</a>.

`))
	})

	t.Run("RedirectWhenError", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(routes.RedirectWhenError{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		operator := routes.RedirectWhenError{}
		req, err := factory.NewRequest(operator.Method(), "/", operator)
		NewWithT(t).Expect(err).To(BeNil())

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).To(Equal(`HTTP/0.0 301 Moved Permanently
Location: /other
X-Meta: service-test@1.0.0/RedirectWhenError
Content-Length: 0

`))
	})

	t.Run("Cookies", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(&routes.Cookie{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		operator := routes.Cookie{}
		req, err := factory.NewRequest(operator.Method(), "/", operator)
		NewWithT(t).Expect(err).To(BeNil())

		cookie := &http.Cookie{
			Name:    "token",
			Value:   "test",
			Expires: time.Now().Add(24 * time.Hour),
		}

		req.AddCookie(cookie)

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).To(Equal(
			`HTTP/0.0 204 No Content
Set-Cookie: ` + cookie.String() + `
X-Meta: service-test@1.0.0/Cookie

`))
	})

	t.Run("ReturnOK", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(routes.DataProvider{}, routes.GetByID{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		reqData := struct {
			routes.DataProvider
			routes.GetByID
		}{
			DataProvider: routes.DataProvider{
				ID: "123456",
			},
			GetByID: routes.GetByID{
				Label: []string{"label"},
			},
		}

		req, err := factory.NewRequestWithContext(
			EnableQueryInBodyForGet(bgctx),
			(routes.GetByID{}).Method(),
			reqData.Path(),
			reqData,
		)
		NewWithT(t).Expect(err).To(BeNil())

		dumped, _ := httputil.DumpRequest(req, true)
		NewWithT(t).Expect(string(dumped)).
			To(Equal("GET /123456 HTTP/1.1\r\nContent-Type: application/x-www-form-urlencoded; param=value\r\n\r\nlabel=label"))

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).
			To(Equal(`HTTP/0.0 200 OK
Content-Type: application/json; charset=utf-8
X-Meta: service-test@1.0.0/GetByID

{"id":"123456","label":"label"}
`))
	})

	t.Run("PostReturnOK", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(routes.Create{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		reqData := routes.Create{
			Data: routes.Data{
				ID:    "123456",
				Label: "123",
			},
		}

		req, err := factory.NewRequest((routes.Create{}).Method(), "/", reqData)
		NewWithT(t).Expect(err).To(BeNil())

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).To(Equal(`HTTP/0.0 201 Created
Content-Type: application/json; charset=utf-8
X-Meta: service-test@1.0.0/Create

{"id":"123456","label":"123"}
`))
	})

	t.Run("PostReturnBadRequest", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(routes.Create{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		reqData := routes.Create{
			Data: routes.Data{
				ID: "123456",
			},
		}

		req, err := factory.NewRequest((routes.Create{}).Method(), "/", reqData)
		NewWithT(t).Expect(err).To(BeNil())

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).To(Equal(`HTTP/0.0 400 Bad Request
Content-Type: application/json; charset=utf-8
X-Meta: service-test@1.0.0/Create

{"key":"badRequest","code":400000000,"msg":"invalid parameters","desc":"","canBeTalk":false,"id":"","sources":["service-test@1.0.0"],"fields":[{"field":"label","msg":"missing required field","in":"body"}]}
`))
	})

	t.Run("ReturnNil", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(routes.DataProvider{}, routes.RemoveByID{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		reqData := routes.DataProvider{
			ID: "123456",
		}

		req, err := factory.NewRequest((routes.RemoveByID{}).Method(), reqData.Path(), reqData)
		NewWithT(t).Expect(err).To(BeNil())

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).To(Equal(`HTTP/0.0 500 Internal Server Error
Content-Type: application/json; charset=utf-8
X-Meta: service-test@1.0.0/RemoveByID
X-Num: 1

{"key":"InternalServerError","code":500999001,"msg":"InternalServerError","desc":"","canBeTalk":false,"id":"","sources":["service-test@1.0.0"],"fields":null}
`))
	})

	t.Run("ReturnAttachment", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(routes.DownloadFile{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		op := routes.DownloadFile{}
		req, err := factory.NewRequest(op.Method(), op.Path(), op)
		NewWithT(t).Expect(err).To(BeNil())

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).To(Equal(`HTTP/0.0 200 OK
Content-Disposition: attachment; filename=text.txt
Content-Type: text/plain
X-Meta: service-test@1.0.0/DownloadFile

123123123`))
	})

	t.Run("ReturnWithProcessError", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(routes.DataProvider{}, routes.UpdateByID{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		reqData := routes.DataProvider{
			ID: "123456",
		}

		req, err := factory.NewRequest((routes.GetByID{}).Method(), reqData.Path(), struct {
			routes.DataProvider
			routes.UpdateByID
		}{
			DataProvider: reqData,
			UpdateByID: routes.UpdateByID{
				Data: routes.Data{
					ID:    "11",
					Label: "11",
				},
			},
		})
		NewWithT(t).Expect(err).To(BeNil())

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).To(Equal(`HTTP/0.0 500 Internal Server Error
Content-Type: application/json; charset=utf-8
X-Meta: service-test@1.0.0/UpdateByID

{"key":"UnknownError","code":500000000,"msg":"UnknownError","desc":"something wrong","canBeTalk":false,"id":"","sources":["service-test@1.0.0"],"fields":null}
`))
	})

	t.Run("ReturnWithValidateError", func(t *testing.T) {
		root := kit.NewRouter(Group("/root"))
		root.Register(kit.NewRouter(routes.DataProvider{}, routes.GetByID{}))

		route := NewHttpRouteMeta(root.Routes()[0])
		hdl := NewRouteHandler(meta, route, factory)

		reqData := routes.DataProvider{
			ID: "10",
		}

		req, err := factory.NewRequest((routes.GetByID{}).Method(), reqData.Path(), reqData)
		NewWithT(t).Expect(err).To(BeNil())

		rw := mock.NewMockResponseWriter()
		hdl.ServeHTTP(rw, req)

		NewWithT(t).Expect(string(rw.MustDumpResponse())).To(Equal(`HTTP/0.0 400 Bad Request
Content-Type: application/json; charset=utf-8
X-Meta: service-test@1.0.0/GetByID

{"key":"badRequest","code":400000000,"msg":"invalid parameters","desc":"","canBeTalk":false,"id":"","sources":["service-test@1.0.0"],"fields":[{"field":"id","msg":"string length should be larger than 6, but got invalid value 2","in":"path"}]}
`))
	})
}
