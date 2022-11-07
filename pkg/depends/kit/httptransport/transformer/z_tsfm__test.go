package transformer_test

import (
	"context"
	"testing"
	_ "unsafe"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
)

var (
	bgctx = context.Background()
)

func TestTransformer(t *testing.T) {
	tfs := transformer.Transformers()
	for _, name := range tfs {
		t.Log(name)
	}
}
