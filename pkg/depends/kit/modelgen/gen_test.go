package modelgen_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/machinefi/w3bstream/pkg/depends/gen/codegen"
	. "github.com/machinefi/w3bstream/pkg/depends/kit/modelgen"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

var (
	g *Generator
	f *codegen.File
	m *Model
)

func init() {
	cwd, _ := os.Getwd()
	dir := filepath.Join(cwd, "./__examples__")
	pkg, _ := pkgx.LoadFrom(dir)

	g = New(pkg)
	g.WithComments = true
	g.WithTableName = true
	g.WithTableInterfaces = true
	g.WithMethods = true
	g.Database = "DB"
	g.StructName = "Org"
	g.Scan()
	g.Output(cwd)

	g = New(pkg)
	g.WithComments = true
	g.WithTableName = true
	g.WithTableInterfaces = true
	g.WithMethods = true
	g.Database = "DB"
	g.StructName = "User"
	g.Scan()
	g.Output(cwd)

	f = codegen.NewFile("example", "mock.go") // mock codegen.File
	m = GetModelByName(g, "User")
	if m == nil {
		panic("should scanned `User` model")
	}
}

func ExampleModel_SnippetTableInstanceAndInit() {
	ss := m.SnippetTableInstanceAndInit(f)
	for _, s := range ss {
		fmt.Println(string(s.Bytes()))
	}
	// Output:
	// var UserTable *builder.Table
	// func init() {
	// UserTable=DB.Register(&User{})
	// }
}

func ExampleModel_SnippetTableIteratorAndMethods() {
	ss := m.SnippetTableIteratorAndMethods(f)
	for _, s := range ss {
		fmt.Println(string(s.Bytes()))
	}
	// Output:
	// type UserIterator struct {
	// }
	// func (*UserIterator) New() interface{} {
	// return &User{}
	// }
	// func (*UserIterator) Resolve(v interface{}) *User {
	// return v.(*User)
	// }
}

func ExampleModel_SnippetTableName() {
	fmt.Println(string(m.SnippetTableName(f).Bytes()))
	// Output:
	// func (*User) TableName() string {
	// return "t_user"
	// }
}

func ExampleModel_SnippetTableDesc() {
	fmt.Println(string(m.SnippetTableDesc(f).Bytes()))
	// Output:
	// func (*User) TableDesc() []string {
	// return []string{
	// "User 用户表",
	// }
	// }
}

func ExampleModel_SnippetComments() {
	fmt.Println(string(m.SnippetComments(f).Bytes()))
	// Output:
	// func (*User) Comments() map[string]string {
	// return map[string]string{
	// "Name": "姓名",
	// "Nickname": "昵称",
	// "OrgID": "关联组织",
	// "Username": "用户名",
	// }
	// }
}

func ExampleModel_SnippetColDesc() {
	fmt.Println(string(m.SnippetColDesc(f).Bytes()))
	// Output:
	// func (*User) ColDesc() map[string][]string {
	// return map[string][]string{
	// "Name": []string{
	// "姓名",
	// },
	// "Nickname": []string{
	// "昵称",
	// },
	// "OrgID": []string{
	// "关联组织",
	// "组织ID",
	// },
	// "Username": []string{
	// "用户名",
	// },
	// }
	// }
}

func ExampleModel_SnippetColRel() {
	fmt.Println(string(m.SnippetColRel(f).Bytes()))
	// Output:
	// func (*User) ColRel() map[string][]string {
	// return map[string][]string{
	// "OrgID": []string{
	// "Org",
	// "ID",
	// },
	// }
	// }
}

func ExampleModel_SnippetPrimaryKey() {
	fmt.Println(string(m.SnippetPrimaryKey(f).Bytes()))
	// Output:
	// func (*User) PrimaryKey() []string {
	// return []string{
	// "ID",
	// }
	// }
}

func ExampleModel_SnippetIndexes() {
	fmt.Println(string(m.SnippetIndexes(f).Bytes()))
	// Output:
	// func (*User) Indexes() builder.Indexes {
	// return builder.Indexes{
	// "i_geom/SPATIAL": []string{
	// "(#Geom)",
	// },
	// "i_nickname/BTREE": []string{
	// "Name",
	// },
	// "i_username": []string{
	// "Username",
	// },
	// }
	// }
}

func ExampleModel_SnippetIndexFieldNames() {
	fmt.Println(string(m.SnippetIndexFieldNames(f).Bytes()))
	// Output:
	// func (m *User) IndexFieldNames() []string {
	// return []string{
	// "ID",
	// "Name",
	// "OrgID",
	// "Username",
	// }
	// }
}

func ExampleModel_SnippetUniqueIndexes() {
	ss := m.SnippetUniqueIndexes(f)
	for _, s := range ss {
		fmt.Println(string(s.Bytes()))
	}
	// Output:
	// func (*User) UniqueIndexes() builder.Indexes {
	// return builder.Indexes{
	// "ui_id_org": []string{
	// "ID",
	// "OrgID",
	// "DeletedAt",
	// },
	// "ui_name": []string{
	// "Name",
	// "DeletedAt",
	// },
	// }
	// }
	// func (*User) UniqueIndexUIIDOrg() string {
	// return "ui_id_org"
	// }
	// func (*User) UniqueIndexUIName() string {
	// return "ui_name"
	// }
}

func ExampleModel_SnippetFieldMethods() {
	ss := m.SnippetFieldMethods(f)
	for _, s := range ss {
		fmt.Print(string(s.Bytes()))
	}
	// Output:
	// func (m *User) ColID() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldID())
	// }func (*User) FieldID() string {
	// return "ID"
	// }func (m *User) ColName() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldName())
	// }func (*User) FieldName() string {
	// return "Name"
	// }func (m *User) ColNickname() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldNickname())
	// }func (*User) FieldNickname() string {
	// return "Nickname"
	// }func (m *User) ColUsername() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldUsername())
	// }func (*User) FieldUsername() string {
	// return "Username"
	// }func (m *User) ColGender() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldGender())
	// }func (*User) FieldGender() string {
	// return "Gender"
	// }func (m *User) ColBoolean() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldBoolean())
	// }func (*User) FieldBoolean() string {
	// return "Boolean"
	// }func (m *User) ColGeom() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldGeom())
	// }func (*User) FieldGeom() string {
	// return "Geom"
	// }func (m *User) ColOrgID() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldOrgID())
	// }func (*User) FieldOrgID() string {
	// return "OrgID"
	// }func (m *User) ColCreatedAt() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldCreatedAt())
	// }func (*User) FieldCreatedAt() string {
	// return "CreatedAt"
	// }func (m *User) ColUpdatedAt() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldUpdatedAt())
	// }func (*User) FieldUpdatedAt() string {
	// return "UpdatedAt"
	// }func (m *User) ColDeletedAt() *builder.Column {
	// return UserTable.ColByFieldName(m.FieldDeletedAt())
	// }func (*User) FieldDeletedAt() string {
	// return "DeletedAt"
	// }
}

func ExampleModel_SnippetCondByValue() {
	fmt.Println(string(m.SnippetCondByValue(f).Bytes()))
	// Output:
	// func (m *User) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
	// var (
	// tbl = db.T(m)
	// fvs = builder.FieldValueFromStructByNoneZero(m)
	// cond = []builder.SqlCondition{tbl.ColByFieldName("DeletedAt").Eq(0)}
	// )
	//
	// for _, fn := range m.IndexFieldNames() {
	// if v, ok := fvs[fn]; ok {
	// cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
	// delete(fvs, fn)
	// }
	// }
	// if len(cond) == 0 {
	// panic(fmt.Errorf("no field for indexes has value"))
	// }
	// for fn, v := range fvs {
	// cond = append(cond, tbl.ColByFieldName(fn).Eq(v))
	// }
	// return builder.And(cond...)
	// }
}

func ExampleModel_SnippetCreate() {
	fmt.Println(string(m.SnippetCreate(f).Bytes()))
	// Output:
	// func (m *User) Create(db sqlx.DBExecutor) error {
	//
	// if m.CreatedAt.IsZero() {
	// m.CreatedAt.Set(time.Now())
	// }
	//
	// if m.UpdatedAt.IsZero() {
	// m.UpdatedAt.Set(time.Now())
	// }
	//
	// _, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	// return err
	// }
}

func ExampleModel_SnippetList() {
	fmt.Println(string(m.SnippetList(f).Bytes()))
	// Output:
	// func (m *User) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]User, error) {
	// var (
	// tbl = db.T(m)
	// lst = make([]User, 0)
	// )
	// cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	// adds = append([]builder.Addition{builder.Where(cond), builder.Comment("User.List")}, adds...)
	// err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	// return lst, err
	// }
}

func ExampleModel_SnippetCount() {
	fmt.Print(string(m.SnippetCount(f).Bytes()))
	// Output:
	// func (m *User) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	// tbl := db.T(m)
	// cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	// adds = append([]builder.Addition{builder.Where(cond), builder.Comment("User.List")}, adds...)
	// err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	// return
	// }
}

func ExampleModel_SnippetCRUDByUniqueKeys() {
	ss := m.SnippetCRUDByUniqueKeys(f, "primary", "ui_name")
	for _, s := range ss {
		fmt.Print(string(s.Bytes()))
	}
	// Output:
	// func (m *User) FetchByID(db sqlx.DBExecutor) error {
	// tbl := db.T(m)
	// err := db.QueryAndScan(
	// builder.Select(nil).
	// From(
	// tbl,
	// builder.Where(
	// builder.And(
	// tbl.ColByFieldName("ID").Eq(m.ID),
	// tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
	// ),
	// ),
	// builder.Comment("User.FetchByID"),
	// ),
	// m,
	// )
	// return err
	// }func (m *User) FetchByName(db sqlx.DBExecutor) error {
	// tbl := db.T(m)
	// err := db.QueryAndScan(
	// builder.Select(nil).
	// From(
	// tbl,
	// builder.Where(
	// builder.And(
	// tbl.ColByFieldName("Name").Eq(m.Name),
	// tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
	// ),
	// ),
	// builder.Comment("User.FetchByName"),
	// ),
	// m,
	// )
	// return err
	// }func (m *User) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {
	//
	// if _, ok := fvs["UpdatedAt"]; !ok {
	// fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	// }
	// tbl := db.T(m)
	// res, err := db.Exec(
	// builder.Update(tbl).
	// Where(
	// builder.And(
	// tbl.ColByFieldName("ID").Eq(m.ID),
	// tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
	// ),
	// builder.Comment("User.UpdateByIDWithFVs"),
	// ).
	// Set(tbl.AssignmentsByFieldValues(fvs)...),
	// )
	// if err != nil {
	// return err
	// }
	// if affected, _ := res.RowsAffected(); affected == 0 {
	// return m.FetchByID(db)
	// }
	// return nil
	// }func (m *User) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	// fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	// return m.UpdateByIDWithFVs(db, fvs)
	// }func (m *User) UpdateByNameWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {
	//
	// if _, ok := fvs["UpdatedAt"]; !ok {
	// fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	// }
	// tbl := db.T(m)
	// res, err := db.Exec(
	// builder.Update(tbl).
	// Where(
	// builder.And(
	// tbl.ColByFieldName("Name").Eq(m.Name),
	// tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
	// ),
	// builder.Comment("User.UpdateByNameWithFVs"),
	// ).
	// Set(tbl.AssignmentsByFieldValues(fvs)...),
	// )
	// if err != nil {
	// return err
	// }
	// if affected, _ := res.RowsAffected(); affected == 0 {
	// return m.FetchByName(db)
	// }
	// return nil
	// }func (m *User) UpdateByName(db sqlx.DBExecutor, zeros ...string) error {
	// fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	// return m.UpdateByNameWithFVs(db, fvs)
	// }func (m *User) Delete(db sqlx.DBExecutor) error {
	// _, err := db.Exec(
	// builder.Delete().
	// From(
	// db.T(m),
	// builder.Where(m.CondByValue(db)),
	// builder.Comment("User.Delete"),
	// ),
	// )
	// return err
	// }func (m *User) DeleteByID(db sqlx.DBExecutor) error {
	// tbl := db.T(m)
	// _, err := db.Exec(
	// builder.Delete().
	// From(
	// tbl,
	// builder.Where(
	// builder.And(
	// tbl.ColByFieldName("ID").Eq(m.ID),
	// tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
	// ),
	// ),
	// builder.Comment("User.DeleteByID"),
	// ),
	// )
	// return err
	// }func (m *User) SoftDeleteByID(db sqlx.DBExecutor) error {
	// tbl := db.T(m)
	// fvs := builder.FieldValues{}
	//
	// if _, ok := fvs["DeletedAt"]; !ok {
	// fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	// }
	//
	// if _, ok := fvs["UpdatedAt"]; !ok {
	// fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	// }
	// _, err := db.Exec(
	// builder.Update(db.T(m)).
	// Where(
	// builder.And(
	// tbl.ColByFieldName("ID").Eq(m.ID),
	// tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
	// ),
	// builder.Comment("User.SoftDeleteByID"),
	// ).
	// Set(tbl.AssignmentsByFieldValues(fvs)...),
	// )
	// return err
	// }func (m *User) DeleteByName(db sqlx.DBExecutor) error {
	// tbl := db.T(m)
	// _, err := db.Exec(
	// builder.Delete().
	// From(
	// tbl,
	// builder.Where(
	// builder.And(
	// tbl.ColByFieldName("Name").Eq(m.Name),
	// tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
	// ),
	// ),
	// builder.Comment("User.DeleteByName"),
	// ),
	// )
	// return err
	// }func (m *User) SoftDeleteByName(db sqlx.DBExecutor) error {
	// tbl := db.T(m)
	// fvs := builder.FieldValues{}
	//
	// if _, ok := fvs["DeletedAt"]; !ok {
	// fvs["DeletedAt"] = types.Timestamp{Time: time.Now()}
	// }
	//
	// if _, ok := fvs["UpdatedAt"]; !ok {
	// fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	// }
	// _, err := db.Exec(
	// builder.Update(db.T(m)).
	// Where(
	// builder.And(
	// tbl.ColByFieldName("Name").Eq(m.Name),
	// tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
	// ),
	// builder.Comment("User.SoftDeleteByName"),
	// ).
	// Set(tbl.AssignmentsByFieldValues(fvs)...),
	// )
	// return err
	// }
}
