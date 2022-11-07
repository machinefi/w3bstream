package builder

import (
	"database/sql/driver"
	"reflect"
)

type Model interface {
	TableName() string
}

type WithPrimaryKey interface {
	PrimaryKey() []string
}

type WithUniqueIndexes interface {
	UniqueIndexes() Indexes
}

type WithIndexes interface {
	Indexes() Indexes
}

type WithComments interface {
	Comments() map[string]string
}

type WithTableDesc interface {
	TableDesc() []string
}

type WithColDesc interface {
	ColDesc() map[string][]string
}

type WithColRel interface {
	ColRel() map[string][]string
}

type Indexes map[string][]string

var (
	RtModel        = reflect.TypeOf((*Model)(nil)).Elem()
	RtDriverValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
)

type Dialect interface {
	DriverName() string
	PrimaryKeyName() string
	IsErrorUnknownDatabase(error) bool
	IsErrorConflict(error) bool
	CreateDatabase(string) SqlExpr
	CreateSchema(string) SqlExpr
	DropDatabase(string) SqlExpr
	CreateTableIsNotExists(t *Table) []SqlExpr
	DropTable(*Table) SqlExpr
	TruncateTable(*Table) SqlExpr
	AddColumn(*Column) SqlExpr
	RenameColumn(*Column, *Column) SqlExpr
	ModifyColumn(*Column, *Column) SqlExpr
	DropColumn(*Column) SqlExpr
	AddIndex(*Key) SqlExpr
	DropIndex(*Key) SqlExpr
	DataType(*ColumnType) SqlExpr
}

type DataTypeDescriber interface {
	DataType(driver string) string
}
