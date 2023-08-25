// This is a generated source file. DO NOT EDIT
// Source: models/proof__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var ProofTable *builder.Table

func init() {
	ProofTable = DB.Register(&Proof{})
}

type ProofIterator struct {
}

func (*ProofIterator) New() interface{} {
	return &Proof{}
}

func (*ProofIterator) Resolve(v interface{}) *Proof {
	return v.(*Proof)
}

func (*Proof) TableName() string {
	return "t_proof"
}

func (*Proof) TableDesc() []string {
	return []string{
		"Proof database model proof",
	}
}

func (*Proof) Comments() map[string]string {
	return map[string]string{}
}

func (*Proof) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*Proof) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*Proof) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *Proof) IndexFieldNames() []string {
	return []string{
		"ID",
		"ProofID",
	}
}

func (*Proof) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_proof_id": []string{
			"ProofID",
		},
	}
}

func (*Proof) UniqueIndexUIProofID() string {
	return "ui_proof_id"
}

func (m *Proof) ColID() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldID())
}

func (*Proof) FieldID() string {
	return "ID"
}

func (m *Proof) ColProjectID() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldProjectID())
}

func (*Proof) FieldProjectID() string {
	return "ProjectID"
}

func (m *Proof) ColProofID() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldProofID())
}

func (*Proof) FieldProofID() string {
	return "ProofID"
}

func (m *Proof) ColName() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldName())
}

func (*Proof) FieldName() string {
	return "Name"
}

func (m *Proof) ColTemplateName() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldTemplateName())
}

func (*Proof) FieldTemplateName() string {
	return "TemplateName"
}

func (m *Proof) ColImageID() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldImageID())
}

func (*Proof) FieldImageID() string {
	return "ImageID"
}

func (m *Proof) ColInputData() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldInputData())
}

func (*Proof) FieldInputData() string {
	return "InputData"
}

func (m *Proof) ColReceipt() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldReceipt())
}

func (*Proof) FieldReceipt() string {
	return "Receipt"
}

func (m *Proof) ColStatus() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldStatus())
}

func (*Proof) FieldStatus() string {
	return "Status"
}

func (m *Proof) ColCreatedAt() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldCreatedAt())
}

func (*Proof) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *Proof) ColUpdatedAt() *builder.Column {
	return ProofTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*Proof) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *Proof) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *Proof) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *Proof) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]Proof, error) {
	var (
		tbl = db.T(m)
		lst = make([]Proof, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Proof.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *Proof) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("Proof.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *Proof) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Proof.FetchByID"),
			),
		m,
	)
	return err
}

func (m *Proof) FetchByProofID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProofID").Eq(m.ProofID),
					),
				),
				builder.Comment("Proof.FetchByProofID"),
			),
		m,
	)
	return err
}

func (m *Proof) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("Proof.UpdateByIDWithFVs"),
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

func (m *Proof) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *Proof) UpdateByProofIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ProofID").Eq(m.ProofID),
				),
				builder.Comment("Proof.UpdateByProofIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByProofID(db)
	}
	return nil
}

func (m *Proof) UpdateByProofID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProofIDWithFVs(db, fvs)
}

func (m *Proof) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("Proof.Delete"),
			),
	)
	return err
}

func (m *Proof) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("Proof.DeleteByID"),
			),
	)
	return err
}

func (m *Proof) DeleteByProofID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProofID").Eq(m.ProofID),
					),
				),
				builder.Comment("Proof.DeleteByProofID"),
			),
	)
	return err
}
