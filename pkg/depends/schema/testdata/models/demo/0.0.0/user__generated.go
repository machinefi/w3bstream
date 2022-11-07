// This is a generated source file. DO NOT EDIT
// Source: demo/user__generated.go

package demo

import "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"

var UserTable *builder.Table

func init() {
	UserTable = DB.Register(&User{})
}

type UserIterator struct {
}

func (UserIterator) New() interface{} {
	return &User{}
}

func (UserIterator) Resolve(v interface{}) *User {
	return v.(*User)
}

func (*User) TableName() string {
	return "t_user"
}

func (*User) TableDesc() []string {
	return []string{
		"User User Table",
	}
}

func (*User) Comments() map[string]string {
	return map[string]string{
		"Gender":   "gender 1 male 2 female",
		"ID":       "user id",
		"Username": "username",
	}
}

func (*User) ColDesc() map[string][]string {
	return map[string][]string{
		"Gender": []string{
			"gender 1 male 2 female",
		},
		"ID": []string{
			"user id",
		},
		"Username": []string{
			"username",
		},
	}
}

func (*User) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*User) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *User) IndexFieldNames() []string {
	return []string{
		"ID",
		"Username",
	}
}

func (*User) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_username": []string{
			"Username",
			"DeletedAt",
		},
	}
}

func (*User) UniqueIndexUIUsername() string {
	return "ui_username"
}

func (m *User) ColID() *builder.Column {
	return UserTable.ColByFieldName(m.FieldID())
}

func (*User) FieldID() string {
	return "ID"
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
