// This is a generated source file. DO NOT EDIT
// Source: models/publisher__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/kit/sqlx"
	"github.com/iotexproject/Bumblebee/kit/sqlx/builder"
)

var PublisherTable *builder.Table

func init() {
	PublisherTable = DB.Register(&Publisher{})
}

type PublisherIterator struct {
}

func (PublisherIterator) New() interface{} {
	return &Publisher{}
}

func (PublisherIterator) Resolve(v interface{}) *Publisher {
	return v.(*Publisher)
}

func (*Publisher) TableName() string {
	return "t_publisher"
}

func (*Publisher) TableDesc() []string {
	return []string{
		"Publisher database model demo",
	}
}

func (*Publisher) Comments() map[string]string {
	return map[string]string{}
}

func (*Publisher) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Publisher) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Publisher) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Publisher) IndexFieldNames() []string {
	return []string{
		"ID",
		"PublisherID",
	}
}

func (*Publisher) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_publisher_id": []string{
			"PublisherID",
		},
	}
}

func (*Publisher) UniqueIndexUIPublisherID() string {
	return "ui_publisher_id"
}

func (m *Publisher) ColID() *builder.Column {
	return PublisherTable.ColByFieldName(m.FieldID())
}

func (*Publisher) FieldID() string {
	return "ID"
}

func (m *Publisher) ColProjectID() *builder.Column {
	return PublisherTable.ColByFieldName(m.FieldProjectID())
}

func (*Publisher) FieldProjectID() string {
	return "ProjectID"
}

func (m *Publisher) ColPublisherID() *builder.Column {
	return PublisherTable.ColByFieldName(m.FieldPublisherID())
}

func (*Publisher) FieldPublisherID() string {
	return "PublisherID"
}

func (m *Publisher) ColName() *builder.Column {
	return PublisherTable.ColByFieldName(m.FieldName())
}

func (*Publisher) FieldName() string {
	return "Name"
}

func (m *Publisher) ColKey() *builder.Column {
	return PublisherTable.ColByFieldName(m.FieldKey())
}

func (*Publisher) FieldKey() string {
	return "Key"
}

func (m *Publisher) ColToken() *builder.Column {
	return PublisherTable.ColByFieldName(m.FieldToken())
}

func (*Publisher) FieldToken() string {
	return "Token"
}

func (m *Publisher) ColCreatedAt() *builder.Column {
	return PublisherTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Publisher) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Publisher) ColUpdatedAt() *builder.Column {
	return PublisherTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Publisher) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Publisher) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
	var (
		tbl  = db.T(m)
		fvs  = builder.FieldValueFromStructByNoneZero(m)
		cond = make([]builder.SqlCondition, 0)
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

func (m *Publisher) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Publisher) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Publisher, error) {
	var (
		tbl = db.T(m)
		lst = make([]Publisher, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Publisher.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Publisher) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Publisher.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Publisher) FetchByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
					),
				),
				builder.Comment("Publisher.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Publisher) FetchByPublisherID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("PublisherID").Eq(m.PublisherID),
					),
				),
				builder.Comment("Publisher.FetchByPublisherID"),
			),
		m,
	)
	return err
}

func (m *Publisher) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ID").Eq(m.ID),
				),
				builder.Comment("Publisher.UpdateByIDWithFVs"),
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

func (m *Publisher) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Publisher) UpdateByPublisherIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("PublisherID").Eq(m.PublisherID),
				),
				builder.Comment("Publisher.UpdateByPublisherIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByPublisherID(db)
	}
	return nil
}

func (m *Publisher) UpdateByPublisherID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByPublisherIDWithFVs(db, fvs)
}

func (m *Publisher) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Publisher.Delete"),
			),
	)
	return err
}

func (m *Publisher) DeleteByID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ID").Eq(m.ID),
					),
				),
				builder.Comment("Publisher.DeleteByID"),
			),
	)
	return err
}

func (m *Publisher) DeleteByPublisherID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("PublisherID").Eq(m.PublisherID),
					),
				),
				builder.Comment("Publisher.DeleteByPublisherID"),
			),
	)
	return err
}
