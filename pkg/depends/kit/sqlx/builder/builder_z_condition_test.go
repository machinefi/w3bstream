package builder_test

import (
	"testing"

	g "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	. "github.com/machinefi/w3bstream/pkg/depends/testutil/buildertestutil"
)

func TestCondition(t *testing.T) {
	t.Run("ChainedConditions", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("a").Eq(1).
				And(nil).
				And(Col("b").LLike("text")).
				Or(Col("a").Eq(2)).
				Xor(Col("b").RLike("g")),
		).To(BeExpr(
			"(((a = ?) AND (b LIKE ?)) OR (a = ?)) XOR (b LIKE ?)",
			1, "%text", 2, "g%",
		))
	})

	t.Run("ComposedCondition", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Xor(
				Or(
					And(
						(*Condition)(nil),
						(*Condition)(nil),
						(*Condition)(nil),
						(*Condition)(nil),
						Col("c").In(1, 2),
						Col("c").In([]int{3, 4}),
						Col("a").Eq(1),
						Col("b").Like("text"),
					),
					Col("a").Eq(2),
				),
				Col("b").Like("g"),
			),
		).To(BeExpr(
			"(((c IN (?,?)) AND (c IN (?,?)) AND (a = ?) AND (b LIKE ?)) OR (a = ?)) XOR (b LIKE ?)",
			1, 2, 3, 4, 1, "%text%", 2, "%g%",
		))
	})
	t.Run("SlipNil", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Xor(
				Col("a").In(),
				Or(
					Col("a").NotIn(),
					And(
						nil,
						Col("a").Eq(1),
						Col("b").Like("text"),
					),
					Col("a").Eq(2),
				),
				Col("b").Like("g"),
			),
		).To(BeExpr(
			"(((a = ?) AND (b LIKE ?)) OR (a = ?)) XOR (b LIKE ?)",
			1, "%text%", 2, "%g%",
		))
	})
	t.Run("XorAndOr", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Xor(
				Col("a").In(),
				Or(
					Col("a").NotIn(),
					And(
						nil,
						Col("a").Eq(1),
						Col("b").Like("text"),
					),
					Col("a").Eq(2),
				),
				Col("b").Like("g"),
			),
		).To(BeExpr(
			"(((a = ?) AND (b LIKE ?)) OR (a = ?)) XOR (b LIKE ?)",
			1, "%text%", 2, "%g%",
		))
	})
	t.Run("XOR", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Xor(
				Col("a").Eq(1),
				Col("b").Like("g"),
			),
		).To(BeExpr(
			"(a = ?) XOR (b LIKE ?)",
			1, "%g%",
		))
	})
	t.Run("Like", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").Like("e"),
		).To(BeExpr(
			"d LIKE ?",
			"%e%",
		))
	})
	t.Run("NotLike", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").NotLike("e"),
		).To(BeExpr(
			"d NOT LIKE ?",
			"%e%",
		))
	})
	t.Run("Equal", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").Eq("e"),
		).To(BeExpr(
			"d = ?",
			"e",
		))
	})
	t.Run("NotEqual", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").Neq("e"),
		).To(BeExpr(
			"d <> ?",
			"e",
		))
	})
	t.Run("In", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").In("e", "f"),
		).To(BeExpr(
			"d IN (?,?)",
			"e", "f",
		))
	})
	t.Run("NotIn", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").NotIn("e", "f"),
		).To(BeExpr(
			"d NOT IN (?,?)",
			"e", "f",
		))
	})
	t.Run("LessThan", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").Lt(3),
		).To(BeExpr(
			"d < ?",
			3,
		))
	})
	t.Run("LessOrEqualThan", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").Lte(3),
		).To(BeExpr(
			"d <= ?",
			3,
		))
	})
	t.Run("GreaterThan", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").Gt(3),
		).To(BeExpr(
			"d > ?",
			3,
		))
	})
	t.Run("GreaterOrEqualThan", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").Gte(3),
		).To(BeExpr(
			"d >= ?",
			3,
		))
	})
	t.Run("Between", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").Between(0, 2),
		).To(BeExpr(
			"d BETWEEN ? AND ?",
			0, 2,
		))
	})
	t.Run("NotBetween", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").NotBetween(0, 2),
		).To(BeExpr(
			"d NOT BETWEEN ? AND ?",
			0, 2,
		))
	})
	t.Run("IsNull", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").IsNull(),
		).To(BeExpr(
			"d IS NULL",
		))
	})
	t.Run("IsNotNull", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Col("d").IsNotNull(),
		).To(BeExpr(
			"d IS NOT NULL",
		))
	})
}
