package swagger

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
)

var raw = bytes.NewBuffer(nil)

func init() {
	data, err := ioutil.ReadFile("./openapi.json")
	if err == nil {
		raw.Write(data)
	} else {
		raw.Write([]byte("{}"))
	}
}

var Router = kit.NewRouter(OpenAPI{})

type OpenAPI struct {
	httpx.MethodGet
}

func (s OpenAPI) Output(ctx context.Context) (interface{}, error) {
	return httpx.WrapContentType(httpx.MIME_JSON)(bytes.NewBuffer(raw.Bytes())), nil
}
