// This is a generated source file. DO NOT EDIT
// Source: models/tag__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var TagTable *builder.Table

func init() {
	TagTable = DB.Register(&Tag{})
}

type TagIterator struct {
}

func (*TagIterator) New() interface{} {
	return &Tag{}
}

func (*TagIterator) Resolve(v interface{}) *Tag {
	return v.(*Tag)
}

func (*Tag) TableName() string {
	return "t_tag"
}

func (*Tag) TableDesc() []string {
	return []string{
		"Tag tag for other object",
	}
}

func (*Tag) Comments() map[string]string {
	return map[string]string{}
}

func (*Tag) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Tag) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Tag) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Tag) IndexFieldNames() []string {
	return []string{
		"ID",
		"TagID",
	}
}

func (*Tag) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_tag_id": []string{
			"TagID",
			"DeletedAt",
		},
	}
}

func (*Tag) UniqueIndexUITagID() string {
	return "ui_tag_id"
}

func (m *Tag) ColID() *builder.Column {
	return TagTable.ColByFieldName(m.FieldID())
}

func (*Tag) FieldID() string {
	return "ID"
}

func (m *Tag) ColProjectID() *builder.Column {
	return TagTable.ColByFieldName(m.FieldProjectID())
}

func (*Tag) FieldProjectID() string {
	return "ProjectID"
}

func (m *Tag) ColTagID() *builder.Column {
	return TagTable.ColByFieldName(m.FieldTagID())
}

func (*Tag) FieldTagID() string {
	return "TagID"
}

func (m *Tag) ColReferenceID() *builder.Column {
	return TagTable.ColByFieldName(m.FieldReferenceID())
}

func (*Tag) FieldReferenceID() string {
	return "ReferenceID"
}

func (m *Tag) ColReferenceType() *builder.Column {
	return TagTable.ColByFieldName(m.FieldReferenceType())
}

func (*Tag) FieldReferenceType() string {
	return "ReferenceType"
}

func (m *Tag) ColInfo() *builder.Column {
	return TagTable.ColByFieldName(m.FieldInfo())
}

func (*Tag) FieldInfo() string {
	return "Info"
}

func (m *Tag) ColCreatedAt() *builder.Column {
	return TagTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Tag) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Tag) ColUpdatedAt() *builder.Column {
	return TagTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Tag) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Tag) ColDeletedAt() *builder.Column {
	return TagTable.ColByFieldName(m.FieldDeletedAt())
}

func (*Tag) FieldDeletedAt() string {
	return "DeletedAt"
}

func (m *Tag) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Tag) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Tag) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Tag, error) {
	var (
		tbl = db.T(m)
		lst = make([]Tag, 0)
	)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Tag.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Tag) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	cond = builder.And(tbl.ColByFieldName("DeletedAt").Eq(0), cond)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Tag.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Tag) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Tag.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Tag) FetchByTagID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("TagID").Eq(m.TagID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Tag.FetchByTagID"),
			),
		m,
	)
	return err
}

func (m *Tag) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Tag.UpdateByIDWithFVs"),
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

func (m *Tag) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Tag) UpdateByTagIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("TagID").Eq(m.TagID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Tag.UpdateByTagIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByTagID(db)
	}
	return nil
}

func (m *Tag) UpdateByTagID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByTagIDWithFVs(db, fvs)
}

func (m *Tag) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Tag.Delete"),
			),
	)
	return err
}

func (m *Tag) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Tag.DeleteByID"),
			),
	)
	return err
}

func (m *Tag) SoftDeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Tag.SoftDeleteByID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}

func (m *Tag) DeleteByTagID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("TagID").Eq(m.TagID),
						tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
					),
				),
				builder.Comment("Tag.DeleteByTagID"),
			),
	)
	return err
}

func (m *Tag) SoftDeleteByTagID(db sqlx.DBExecutor) error {
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
					tbl.ColByFieldName("TagID").Eq(m.TagID),
					tbl.ColByFieldName("DeletedAt").Eq(m.DeletedAt),
				),
				builder.Comment("Tag.SoftDeleteByTagID"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	return err
}
