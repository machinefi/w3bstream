package builder_test

import (
	"testing"

	g "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	. "github.com/machinefi/w3bstream/pkg/depends/testutil/buildertestutil"
)

func TestGroupBy(t *testing.T) {
	tbl := T("T")
	cola := Col("f_a")
	colb := Col("F_b")

	t.Run("SelectGroupBy", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Select(nil).From(
				tbl,
				Where(Col("f_A").Eq(1)),
				GroupBy(Col("F_a")).
					Having(Col("F_A").Eq(1)),
			),
		).To(BeExpr(
			`
SELECT * FROM T
WHERE f_a = ?
GROUP BY f_a HAVING f_a = ?
`,
			1, 1,
		))
	})

	t.Run("SelectMultiGroupByAndOrder", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Select(nil).
				From(
					tbl,
					Where(cola.Eq(1)),
					GroupBy(AscOrder(cola), DescOrder(colb)),
				),
		).To(BeExpr(
			`
SELECT * FROM T
WHERE f_a = ?
GROUP BY (f_a) ASC,(f_b) DESC
`,
			1,
		))
	})
}

func TestJoin(t *testing.T) {
	tUser := T("t_user",
		Col("f_id").Type(uint64(0), ",autoincrement"),
		Col("f_name").Type("", ",size=128,default=''"),
		Col("f_org_id").Type("", ",size=128,default=''"),
	)
	tOrg := T("t_org",
		Col("f_id").Type(uint64(0), ",autoincrement"),
		Col("f_name").Type("", ",size=128,default=''"),
	)

	t.Run("JoinOn", func(t *testing.T) {
		expr := Select(
			MultiWith(
				", ",
				Alias(tUser.Col("f_id"), "f_user_id"),
				Alias(tUser.Col("f_name"), "f_user_name"),
				Alias(tOrg.Col("f_id"), "f_org_id"),
				Alias(tOrg.Col("f_name"), "f_org_name"),
			),
		).From(
			tUser,
			Join(Alias(tOrg, "t_org")).
				On(tUser.Col("f_org_id").Eq(tOrg.Col("f_id"))),
		)
		// t.Log(expr.Ex(context.Background()).Query())
		g.NewWithT(t).Expect(expr).To(BeExpr(`
SELECT t_user.f_id AS f_user_id, t_user.f_name AS f_user_name, t_org.f_id AS f_org_id, t_org.f_name AS f_org_name FROM t_user
JOIN t_org AS t_org ON t_user.f_org_id = t_org.f_id
`))
	})

	t.Run("JoinUsing", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Select(nil).
				From(
					tUser,
					Join(tOrg).Using(tUser.Col("f_org_id")),
				),
		).To(BeExpr(`
SELECT * FROM t_user
JOIN t_org USING (f_org_id)
`,
		))
	})
}

func TestLimit(t *testing.T) {
	tbl := T("T")
	t.Run("SelectLimit", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Select(nil).
				From(
					tbl,
					Where(Col("F_a").Eq(1)),
					Limit(1),
				),
		).To(BeExpr(`
SELECT * FROM T
WHERE f_a = ?
LIMIT 1
`, 1,
		))
	})
	t.Run("SelectWithoutLimit", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Select(nil).
				From(
					tbl,
					Where(Col("F_a").Eq(1)),
					Limit(-1),
				),
		).To(BeExpr(`
SELECT * FROM T
WHERE f_a = ?
`, 1,
		))
	})
	t.Run("SelectLimitAndOffset", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Select(nil).
				From(
					tbl,
					Where(Col("F_a").Eq(1)),
					Limit(1).Offset(200),
				),
		).To(BeExpr(`
SELECT * FROM T
WHERE f_a = ?
LIMIT 1 OFFSET 200
`, 1,
		))
	})
}

func TestOrderBy(t *testing.T) {
	tbl := T("T")
	t.Run("SelectOrder", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Select(nil).
				From(
					tbl,
					OrderBy(
						AscOrder(Col("F_a")),
						DescOrder(Col("F_b")),
					),
					Where(Col("F_a").Eq(1)),
				),
		).To(BeExpr(`
SELECT * FROM T
WHERE f_a = ?
ORDER BY (f_a) ASC,(f_b) DESC
`, 1,
		))
	})
}
