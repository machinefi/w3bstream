package nullable_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/scanner/nullable"
)

func BenchmarkNewNullIgnoreScanner(b *testing.B) {
	v := 0
	for i := 0; i < b.N; i++ {
		_ = nullable.NewNullIgnoreScanner(&v).Scan(2)
	}
	b.Log(v)
}

func TestNullIgnoreScanner(t *testing.T) {
	cases := []*struct {
		name string
		src  interface{}
		dst  interface{}
	}{
		{"ScanVal", int(1), 1},
		{"ScanNil", nil, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := 0
			_ = nullable.NewNullIgnoreScanner(&v).Scan(c.src)
			NewWithT(t).Expect(v).To(Equal(c.dst))
		})
	}
}
