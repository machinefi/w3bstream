package builder_test

import (
	"context"
	"testing"

	g "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

func TestFieldValueFromStruct(t *testing.T) {
	type User struct {
		ID       uint64 `db:"f_id"`
		Name     string `db:"f_name"`
		Username string `db:"f_username"`
	}

	user := User{ID: 12345}
	t.Run("FieldValueFromStruct", func(t *testing.T) {
		g.NewWithT(t).Expect(FieldValueFromStruct(user, []string{})).To(g.HaveLen(0))
		values := FieldValueFromStruct(user, []string{"ID"})
		g.NewWithT(t).Expect(values).To(g.Equal(FieldValues{"ID": user.ID}))
	})

	t.Run("FieldValueFromStructByNoneZero", func(t *testing.T) {
		g.NewWithT(t).Expect(FieldValueFromStructByNoneZero(user)).
			To(g.Equal(FieldValues{"ID": user.ID}))

		g.NewWithT(t).Expect(FieldValueFromStructByNoneZero(user, "Username")).
			To(g.Equal(FieldValues{"ID": user.ID, "Username": user.Username}))
	})
}

func TestGetColumnName(t *testing.T) {
	g.NewWithT(t).Expect(GetColumnName("Text", "")).To(g.Equal("f_text"))
	g.NewWithT(t).Expect(GetColumnName("Text", ",size=256")).To(g.Equal("f_text"))
	g.NewWithT(t).Expect(GetColumnName("Text", "f_xxx")).To(g.Equal("f_xxx"))
	g.NewWithT(t).Expect(GetColumnName("Text", "f_xxx,default=''")).To(g.Equal("f_xxx"))
}

func TestParseIndexDefine(t *testing.T) {
	t.Run("IdxWithFieldNames", func(t *testing.T) {
		i := ParseIndexDefine("index i_xxx/BTREE Name")
		g.NewWithT(t).Expect(i).To(g.Equal(&IndexDefine{
			Kind:     "index",
			Name:     "i_xxx",
			Method:   "BTREE",
			IndexDef: IndexDef{FieldNames: []string{"Name"}},
		}))
	})
	t.Run("PrimaryWithFieldNames", func(t *testing.T) {
		i := ParseIndexDefine("primary ID Name")
		g.NewWithT(t).Expect(i).To(g.Equal(&IndexDefine{
			Kind:     "primary",
			Name:     "",
			Method:   "",
			IndexDef: IndexDef{FieldNames: []string{"ID", "Name"}},
		}))
	})
	t.Run("IndexWithExpr", func(t *testing.T) {
		i := ParseIndexDefine("index i_xxx USING GIST (#TEST gist_trgm_ops)")
		g.NewWithT(t).Expect(i).To(g.Equal(&IndexDefine{
			Kind:     "index",
			Name:     "i_xxx",
			Method:   "",
			IndexDef: IndexDef{Expr: "USING GIST (#TEST gist_trgm_ops)"},
		}))
	})
}

type User struct {
	ID       uint64 `db:"f_id"`
	Name     string `db:"f_name"`
	Username string `db:"f_username"`
}

func (User) TableName() string { return "t_user" }

type OrgUser struct {
	OrgID  uint64 `db:"f_org_id"`
	UserID uint64 `db:"f_user_id"`
}

func (OrgUser) TableName() string { return "t_org_user" }

type Org struct {
	ID   uint64 `db:"f_id"`
	Name string `db:"f_name"`
}

func (Org) TableName() string { return "t_org" }

type OrgUserAll struct {
	OrgUser
	User User `json:"user"`
	Org  Org  `json:"org"`
}

func TestColumnsByStruct(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		q := ColumnsByStruct(&User{}).Ex(context.Background()).Query()
		g.NewWithT(t).Expect(q).
			To(g.Equal("t_user.f_id AS t_user__f_id, t_user.f_name AS t_user__f_name, t_user.f_username AS t_user__f_username"))
	})
	t.Run("Joined", func(t *testing.T) {
		q := ColumnsByStruct(&OrgUserAll{}).Ex(context.Background()).Query()
		// for _, sub := range strings.Split(q, ", ") {
		// 	t.Log(sub)
		// }
		// t.Log(q)
		g.NewWithT(t).Expect(q).To(
			g.Equal("t_org_user.f_org_id AS t_org_user__f_org_id, " +
				"t_org_user.f_user_id AS t_org_user__f_user_id, " +
				"t_user.f_id AS t_user__f_id, " +
				"t_user.f_name AS t_user__f_name, " +
				"t_user.f_username AS t_user__f_username, " +
				"t_org.f_id AS t_org__f_id, " +
				"t_org.f_name AS t_org__f_name",
			),
		)
	})
}
