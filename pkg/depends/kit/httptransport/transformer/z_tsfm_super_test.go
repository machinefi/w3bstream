package transformer_test

import (
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
)

func BenchmarkTransformerSuper(b *testing.B) {
	ts := NewSuper(&PlainText{}, &CommonOption{Omitempty: true})

	b.Run("DecodeFromBySuper", func(b *testing.B) {
		ret := ""
		for i := 0; i < b.N; i++ {
			_ = ts.DecodeFrom(bgctx, NewStringReader("111"), &ret)
		}
		// b.Log(ret)
	})

	b.Run("DecodeFromDirect", func(b *testing.B) {
		pt := PlainText{}
		ret := ""
		for i := 0; i < b.N; i++ {
			_ = pt.DecodeFrom(bgctx, NewStringReader("111"), &ret)
		}
		// b.Log(ret)
	})
}
