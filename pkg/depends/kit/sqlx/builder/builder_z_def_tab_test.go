package builder_test

import (
	"testing"

	g "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	. "github.com/machinefi/w3bstream/pkg/depends/testutil/buildertestutil"
)

func TestTable_Expr(t *testing.T) {
	tUser := T("t_user",
		Col("f_id").Field("ID").Type(uint64(0), ",autoincrement"),
		Col("f_name").Field("Name").Type("", ",size=128,default=''"),
	)

	tUserRole := T("t_user_role",
		Col("f_id").Field("ID").Type(uint64(0), ",autoincrement"),
		Col("f_user_id").Field("UserID").Type(uint64(0), ""),
	)

	t.Run("ReplaceTable", func(t *testing.T) {
		g.NewWithT(t).Expect(tUser.Expr("#.*")).
			To(BeExpr("t_user.*"))
	})
	t.Run("ReplaceTableColByField", func(t *testing.T) {
		g.NewWithT(t).Expect(tUser.Expr("#ID = #ID + 1")).
			To(BeExpr("f_id = f_id + 1"))
	})
	t.Run("ReplaceTableColByFieldForFn", func(t *testing.T) {
		g.NewWithT(t).Expect(tUser.Expr("COUNT(#ID)")).
			To(BeExpr("COUNT(f_id)"))
	})
	t.Run("CouldHandleContext", func(t *testing.T) {
		g.NewWithT(t).Expect(
			Select(nil).
				From(
					tUser,
					Where(AsCond(tUser.Expr("#ID > 1"))),
					Join(tUserRole).
						On(AsCond(
							tUser.Expr("#ID = ?", tUserRole.Expr("#UserID")),
						)),
				),
		).To(BeExpr(`
SELECT * FROM t_user
JOIN t_user_role ON t_user.f_id = t_user_role.f_user_id
WHERE t_user.f_id > 1
`,
		))
	})

	// TODO table diff
}
