package builder_test

import (
	"testing"

	"github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	. "github.com/machinefi/w3bstream/pkg/depends/testutil/buildertestutil"
)

func TestStmtDelete(t *testing.T) {
	table := T("T")

	t.Run("delete", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Delete().From(table,
				Where(Col("F_a").Eq(1)),
				Comment("Comment"),
			),
		).To(BeExpr(`
DELETE FROM T
WHERE f_a = ?
/* Comment */
`, 1))
	})
}

func TestStmtInsert(t *testing.T) {
	table := T("T", Col("f_a"), Col("f_b"))

	t.Run("Simple", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Insert().
				Into(table, Comment("Comment")).
				Values(Cols("f_a", "f_b"), 1, 2),
		).To(BeExpr(`
INSERT INTO T (f_a,f_b) VALUES (?,?)
/* Comment */
`, 1, 2))
	})

	t.Run("WithModifier", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Insert("IGNORE").
				Into(table).
				Values(Cols("f_a", "f_b"), 1, 2),
		).To(BeExpr("INSERT IGNORE INTO T (f_a,f_b) VALUES (?,?)",
			1, 2))
	})

	t.Run("MultiInsert", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Insert().
				Into(table).
				Values(Cols("f_a", "f_b"), 1, 2, 1, 2, 1, 2),
		).To(BeExpr("INSERT INTO T (f_a,f_b) VALUES (?,?),(?,?),(?,?)", 1, 2, 1, 2, 1, 2))
	})

	t.Run("FromSelect", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Insert().
				Into(table).
				Values(Cols("f_a", "f_b"), Select(Cols("f_a", "f_b")).From(table, Where(table.Col("f_a").Eq(1)))),
		).To(BeExpr(`
INSERT INTO T (f_a,f_b) SELECT f_a,f_b FROM T
WHERE f_a = ?
`, 1))
	})
}

func TestSelect(t *testing.T) {
	table := T("T")

	t.Run("Simple", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Select(nil).
				From(
					table,
					Where(
						Col("F_a").Eq(1),
					),
					Comment("comment"),
				),
		).To(BeExpr(`
SELECT * FROM T
WHERE f_a = ?
/* comment */
`, 1))
	})
	t.Run("WithModifier", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Select(nil, "DISTINCT").
				From(
					table,
					Where(
						Col("F_a").Eq(1),
					),
				),
		).To(BeExpr(`
SELECT DISTINCT * FROM T
WHERE f_a = ?`, 1))
	})
	t.Run("WithColumn", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Select(Col("F_a")).
				From(table,
					Where(
						Col("F_a").Eq(1),
					),
				),
		).To(BeExpr(`
SELECT f_a FROM T
WHERE f_a = ?`, 1))
	})
	t.Run("ForUpdate", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Select(nil).From(
				table,
				Where(Col("F_a").Eq(1)),
				ForUpdate(),
			),
		).To(BeExpr(
			`
SELECT * FROM T
WHERE f_a = ?
FOR UPDATE
`,
			1,
		))
	})
}

func TestStmtUpdate(t *testing.T) {
	table := T("T")

	t.Run("Update", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Update(table).
				Set(
					Col("F_a").ValueBy(1),
					Col("F_b").ValueBy(2),
				).
				Where(
					Col("F_a").Eq(1),
					Comment("Comment"),
				),
		).To(BeExpr(`
UPDATE T SET f_a = ?, f_b = ?
WHERE f_a = ?
/* Comment */`, 1, 2, 1))
	})
}
