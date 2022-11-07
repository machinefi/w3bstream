// This is a generated source file. DO NOT EDIT
// Source: example/user__generated.go

package example

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var UserTable *builder.Table

func init() {
	UserTable = DB.Register(&User{})
}

type UserIterator struct {
}

func (*UserIterator) New() interface{} {
	return &User{}
}

func (*UserIterator) Resolve(v interface{}) *User {
	return v.(*User)
}

func (*User) TableName() string {
	return "t_user"
}

func (*User) TableDesc() []string {
	return []string{
		"User 用户表",
	}
}

func (*User) Comments() map[string]string {
	return map[string]string{
		"Name":     "姓名",
		"Nickname": "昵称",
		"OrgID":    "关联组织",
		"Username": "用户名",
	}
}

func (*User) ColDesc() map[string][]string {
	return map[string][]string{
		"Name": []string{
			"姓名",
		},
		"Nickname": []string{
			"昵称",
		},
		"OrgID": []string{
			"关联组织",
			"组织ID",
		},
		"Username": []string{
			"用户名",
		},
	}
}

func (*User) ColRel() map[string][]string {
	return map[string][]string{
		"OrgID": []string{
			"Org",
			"ID",
		},
	}
}

func (*User) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (*User) Indexes() builder.Indexes {
	return builder.Indexes{
		"i_geom/SPATIAL": []string{
			"(#Geom)",
		},
		"i_nickname/BTREE": []string{
			"Name",
		},
		"i_username": []string{
			"Username",
		},
	}
}

func (m *User) IndexFieldNames() []string {
	return []string{
		"ID",
		"Name",
		"OrgID",
		"Username",
	}
}

func (*User) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_id_org": []string{
			"ID",
			"OrgID",
			"DeletedAt",
		},
		"ui_name": []string{
			"Name",
			"DeletedAt",
		},
	}
}

func (*User) UniqueIndexUIIDOrg() string {
	return "ui_id_org"
}

func (*User) UniqueIndexUIName() string {
	return "ui_name"
}

func (m *User) ColID() *builder.Column {
	return UserTable.ColByFieldName(m.FieldID())
}

func (*User) FieldID() string {
	return "ID"
}

func (m *User) ColName() *builder.Column {
	return UserTable.ColByFieldName(m.FieldName())
}

func (*User) FieldName() string {
	return "Name"
}

func (m *User) ColNickname() *builder.Column {
	return UserTable.ColByFieldName(m.FieldNickname())
}

func (*User) FieldNickname() string {
	return "Nickname"
}

func (m *User) ColUsername() *builder.Column {
	return UserTable.ColByFieldName(m.FieldUsername())
}

func (*User) FieldUsername() string {
	return "Username"
}

func (m *User) ColGender() *builder.Column {
	return UserTable.ColByFieldName(m.FieldGender())
}

func (*User) FieldGender() string {
	return "Gender"
}

func (m *User) ColBoolean() *builder.Column {
	return UserTable.ColByFieldName(m.FieldBoolean())
}

func (*User) FieldBoolean() string {
	return "Boolean"
}

func (m *User) ColGeom() *builder.Column {
	return UserTable.ColByFieldName(m.FieldGeom())
}

func (*User) FieldGeom() string {
	return "Geom"
}

func (m *User) ColOrgID() *builder.Column {
	return UserTable.ColByFieldName(m.FieldOrgID())
}

func (*User) FieldOrgID() string {
	return "OrgID"
}

func (m *User) ColCreatedAt() *builder.Column {
	return UserTable.ColByFieldName(m.FieldCreatedAt())
}

func (*User) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *User) ColUpdatedAt() *builder.Column {
	return UserTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*User) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *User) ColDeletedAt() *builder.Column {
	return UserTable.ColByFieldName(m.FieldDeletedAt())
}

func (*User) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *User) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
	var (
		tbl  = db.T(m)
		fvs  = builder.FieldValueFromStructByNoneZero(m)
		cond = []builder.SqlCondition{tbl.ColByFieldName("DeletedAt").Eq(0)}
	)

	for _, fn := range m.IndexFieldNames() {
		if v, ok := fvs[fn]; ok {
			cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
			delete(fvs, fn)
		}
	}
	if len(cond) == 0 {
		panic(fmt.Errorf("no field for indexes has value"))
	}
	for fn, v := range fvs {
		cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
	}
	return builder.And(cond...)
}

func (m *User) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *User) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]User, error) {
	var (
		tbl = db.T(m)
		lst = make([]User, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("User.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *User) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("User.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *User) FetchByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("User.FetchByID"),
			),
		m,
	)
	return err
}

func (m *User) FetchByIDAndOrgID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
						tbl.ColByFieldName("OrgID").Eq(m.OrgID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("User.FetchByIDAndOrgID"),
			),
		m,
	)
	return err
}

func (m *User) FetchByName(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Name").Eq(m.Name),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("User.FetchByName"),
			),
		m,
	)
	return err
}

func (m *User) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("User.UpdateByIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByID(db)
	}
	return nil
}

func (m *User) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *User) UpdateByIDAndOrgIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
					tbl.ColByFieldName("OrgID").Eq(m.OrgID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("User.UpdateByIDAndOrgIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByIDAndOrgID(db)
	}
	return nil
}

func (m *User) UpdateByIDAndOrgID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDAndOrgIDWithFVs(db, fvs)
}

func (m *User) UpdateByNameWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("Name").Eq(m.Name),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("User.UpdateByNameWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByName(db)
	}
	return nil
}

func (m *User) UpdateByName(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByNameWithFVs(db, fvs)
}

func (m *User) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("User.Delete"),
			),
	)
	return err
}

func (m *User) DeleteByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("User.DeleteByID"),
			),
	)
	return err
}

func (m *User) SoftDeleteByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	fvs := builder.FieldValues{}

	if _, ok := fvs["DeletedAt"]; !ok {
		fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	}

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	_, err := db.Exec(
		builder.Update(db.T(m)).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("User.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *User) DeleteByIDAndOrgID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
						tbl.ColByFieldName("OrgID").Eq(m.OrgID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("User.DeleteByIDAndOrgID"),
			),
	)
	return err
}

func (m *User) SoftDeleteByIDAndOrgID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	fvs := builder.FieldValues{}

	if _, ok := fvs["DeletedAt"]; !ok {
		fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	}

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	_, err := db.Exec(
		builder.Update(db.T(m)).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
					tbl.ColByFieldName("OrgID").Eq(m.OrgID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("User.SoftDeleteByIDAndOrgID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *User) DeleteByName(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("Name").Eq(m.Name),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("User.DeleteByName"),
			),
	)
	return err
}

func (m *User) SoftDeleteByName(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	fvs := builder.FieldValues{}

	if _, ok := fvs["DeletedAt"]; !ok {
		fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	}

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	_, err := db.Exec(
		builder.Update(db.T(m)).
			Where(
				builder.And(
					tbl.ColByFieldName("Name").Eq(m.Name),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("User.SoftDeleteByName"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
