package builder_test

import (
	"context"
	"testing"

	g "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	. "github.com/machinefi/w3bstream/pkg/depends/testutil/buildertestutil"
)

func TestResolveExpr(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		g.NewWithT(t).Expect(ResolveExpr(nil)).To(g.BeNil())
	})
}

type Byte uint8

func TestEx(t *testing.T) {
	t.Run("EmptyQuery", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Expr(""),
		).To(BeExpr(""))
	})

	t.Run("FlattenSlice", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Expr(`#ID IN (?)`, []int{28, 29, 30}),
		).To(BeExpr("#ID IN (?,?,?)", 28, 29, 30))
	})

	t.Run("FlattenSliceWithNamedByte", func(t *testing.T) {
		g.NewWithT(t).Expect(
			And(
				And(nil, Col("f_id").In([]int{28})),
				Col("f_id").In([]Byte{28}),
			),
		).To(BeExpr("((f_id IN (?))) AND (f_id IN (?))", 28, Byte(28)))
	})

	t.Run("FlattenSliceShouldSkipBytes", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Expr(`#ID = (?)`, []byte("")),
		).To(BeExpr("#ID = (?)", []byte("")))
	})

	t.Run("FlattenWithSubExpr ", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Expr(`#ID = ?`, Expr("#ID + ?", 1)),
		).To(BeExpr("#ID = #ID + ?", 1))
	})

	t.Run("FlattenWithValuerExpr", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Expr(`#Point = ?`, Point{X: 1, Y: 1}),
		).To(BeExpr("#Point = ST_GeomFromText(?)", Point{X: 1, Y: 1}))
	})
}

func BenchmarkEx(b *testing.B) {
	b.Run("EmptyQuery", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Expr("").Ex(context.Background())
		}
	})

	b.Run("FlattenSlice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Expr(`#ID IN (?)`, []int{28, 29, 30}).Ex(context.Background())
		}
	})

	b.Run("FlattenWithSubExpr", func(b *testing.B) {
		b.Run("Raw", func(b *testing.B) {
			eb := Expr("")
			eb.Grow(2)

			eb.WriteQuery("#ID > ?")
			eb.WriteQuery(" AND ")
			eb.WriteQuery("#ID < ?")

			eb.AppendArgs(1, 10)

			rawBuild := func() *Ex {
				return eb.Ex(context.Background())
			}

			clone := func(ex *Ex) *Ex {
				return Expr(ex.Query(), ex.Args()...).Ex(context.Background())
			}

			b.Run("clone", func(b *testing.B) {
				ex := rawBuild()

				for i := 0; i < b.N; i++ {
					_ = clone(ex)
				}
			})
		})

		b.Run("IsNilExpr", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IsNilExpr(Expr(`#ID > ?`, 1))
			}
		})

		b.Run("ByChain", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				e := AsCond(Expr(`#ID > ?`, 1)).And(AsCond(Expr(`#ID < ?`, 10)))
				e.Ex(context.Background())
			}
		})

		b.Run("ByExpr", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				e := And(
					Col("f_id").Lt(0),
					Col("f_id").In([]int{1, 2, 3}),
				)
				e.Ex(context.Background())
			}
		})

		b.Run("ByExprWithoutRecreated", func(b *testing.B) {
			l := Col("f_id").Lt(0)
			r := Col("f_id").In([]int{1, 2, 3})

			b.Run("Single", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					l.Ex(context.Background())
				}
			})

			b.Run("Composed", func(b *testing.B) {
				e := And(l, l, r, r)

				// b.Log(e.Ex(context.Background()).Query())

				for i := 0; i < b.N; i++ {
					e.Ex(context.Background())
				}
			})
		})
	})
}
