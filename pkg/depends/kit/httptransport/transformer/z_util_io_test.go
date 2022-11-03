package transformer_test

import (
	"io"
	"strings"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
)

func BenchmarkBuffers(b *testing.B) {
	b.Run("StringReaders", func(b *testing.B) {
		inputs := strings.Split(strings.Repeat("1", 10), "")

		for i := 0; i < b.N; i++ {
			buffers := NewStringReaders(inputs)
			for i := 0; i < buffers.Len(); i++ {
				r := buffers.NextReader()
				_, _ = io.ReadAll(r)
			}
		}
	})

	b.Run("StringBuilders", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			builders := NewStringBuilders()
			builders.SetN(10)

			for i := 0; i < 10; i++ {
				w := builders.NextWriter()
				_, _ = w.Write([]byte("1"))
			}

			_ = builders.StringSlice()
		}
	})
}

func TestNewBuffers(t *testing.T) {
	inputs := strings.Split(strings.Repeat("1", 10), "")

	t.Run("StringReaders", func(t *testing.T) {
		buffers := NewStringReaders(inputs)
		results := make([]string, 0)
		for i := 0; i < buffers.Len(); i++ {
			r := buffers.NextReader()
			data, _ := io.ReadAll(r)
			results = append(results, string(data))
		}
		NewWithT(t).Expect(results).To(Equal(inputs))
	})

	t.Run("StringBuilders", func(t *testing.T) {
		builders := NewStringBuilders()
		builders.SetN(10)
		for i := 0; i < 10; i++ {
			w := builders.NextWriter()
			_, _ = w.Write([]byte("1"))
		}
		NewWithT(t).Expect(builders.StringSlice()).To(Equal(inputs))
	})
}
