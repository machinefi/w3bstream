package metax_test

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

func TestParseMeta(t *testing.T) {
	t.Run("parse id", func(t *testing.T) {
		meta := ParseMeta("xxxxxx")
		NewWithT(t).Expect(meta.String()).To(Equal("_id=xxxxxx"))
	})

	t.Run("parse meta", func(t *testing.T) {
		meta := ParseMeta("operator=1&operator=2&_id=xxx")
		NewWithT(t).Expect(meta.Get("operator")).To(Equal("1"))
		NewWithT(t).Expect(meta.String()).To(Equal("_id=xxx&operator=1&operator=2"))
	})
}

func TestMeta(t *testing.T) {

	t.Run("ContextConcat", func(t *testing.T) {
		ctx := ContextWith(context.Background(), "key", "1")
		ctx = ContextWithMeta(ctx, (Meta{}).With("key", "2", "3"))

		NewWithT(t).Expect(GetMetaFrom(ctx)["key"]).To(Equal([]string{"1", "2", "3"}))
	})

	t.Run("ContextOverwrite", func(t *testing.T) {
		ctx := ContextWith(context.Background(), "_key", "1")
		ctx = ContextWithMeta(ctx, (Meta{}).With("_key", "2", "3"))

		NewWithT(t).Expect(GetMetaFrom(ctx)["_key"]).To(Equal([]string{"2", "3"}))
	})

	t.Run("EmptyKeyIgnore", func(t *testing.T) {
		ctx := ContextWith(context.Background(), "", "1")
		NewWithT(t).Expect(GetMetaFrom(ctx)).To(HaveLen(0))
	})
}

type Some struct{ Ctx }

func (s *Some) WithContext(ctx context.Context) *Some {
	return &Some{Ctx: s.Ctx.WithContext(ctx)}
}

func TestCtx(t *testing.T) {
	s := &Some{}
	s2 := s.WithContext(ContextWith(context.Background(), "k", "1"))

	NewWithT(t).Expect(GetMetaFrom(s.Context()).Get("k")).To(Equal(""))
	NewWithT(t).Expect(GetMetaFrom(s2.Context()).Get("k")).To(Equal("1"))
}
