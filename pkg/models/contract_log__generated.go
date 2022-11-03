// This is a generated source file. DO NOT EDIT
// Source: models/contract_log__generated.go

package models

import (
	"fmt"
	"time"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

var ContractLogTable *builder.Table

func init() {
	ContractLogTable = MonitorDB.Register(&ContractLog{})
}

type ContractLogIterator struct {
}

func (*ContractLogIterator) New() interface{} {
	return &ContractLog{}
}

func (*ContractLogIterator) Resolve(v interface{}) *ContractLog {
	return v.(*ContractLog)
}

func (*ContractLog) TableName() string {
	return "t_contract_log"
}

func (*ContractLog) TableDesc() []string {
	return []string{
		"ContractLog database model contract log",
	}
}

func (*ContractLog) Comments() map[string]string {
	return map[string]string{}
}

func (*ContractLog) ColDesc() map[string][]string {
	return map[string][]string{}
}

func (*ContractLog) ColRel() map[string][]string {
	return map[string][]string{}
}

func (*ContractLog) PrimaryKey() []string {
	return []string{
		"ID",
	}
}

func (m *ContractLog) IndexFieldNames() []string {
	return []string{
		"ChainID",
		"ContractAddress",
		"ContractLogID",
		"EventType",
		"ID",
		"ProjectName",
		"Topic0",
		"Topic1",
		"Topic2",
		"Topic3",
		"Uniq",
	}
}

func (*ContractLog) UniqueIndexes() builder.Indexes {
	return builder.Indexes{
		"ui_contract_log_id": []string{
			"ContractLogID",
		},
		"ui_contract_log_uniq": []string{
			"ProjectName",
			"EventType",
			"ChainID",
			"ContractAddress",
			"Topic0",
			"Topic1",
			"Topic2",
			"Topic3",
			"Uniq",
		},
	}
}

func (*ContractLog) UniqueIndexUIContractLogID() string {
	return "ui_contract_log_id"
}

func (*ContractLog) UniqueIndexUIContractLogUniq() string {
	return "ui_contract_log_uniq"
}

func (m *ContractLog) ColID() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldID())
}

func (*ContractLog) FieldID() string {
	return "ID"
}

func (m *ContractLog) ColContractLogID() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldContractLogID())
}

func (*ContractLog) FieldContractLogID() string {
	return "ContractLogID"
}

func (m *ContractLog) ColProjectName() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldProjectName())
}

func (*ContractLog) FieldProjectName() string {
	return "ProjectName"
}

func (m *ContractLog) ColUniq() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldUniq())
}

func (*ContractLog) FieldUniq() string {
	return "Uniq"
}

func (m *ContractLog) ColEventType() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldEventType())
}

func (*ContractLog) FieldEventType() string {
	return "EventType"
}

func (m *ContractLog) ColChainID() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldChainID())
}

func (*ContractLog) FieldChainID() string {
	return "ChainID"
}

func (m *ContractLog) ColContractAddress() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldContractAddress())
}

func (*ContractLog) FieldContractAddress() string {
	return "ContractAddress"
}

func (m *ContractLog) ColBlockStart() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldBlockStart())
}

func (*ContractLog) FieldBlockStart() string {
	return "BlockStart"
}

func (m *ContractLog) ColBlockCurrent() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldBlockCurrent())
}

func (*ContractLog) FieldBlockCurrent() string {
	return "BlockCurrent"
}

func (m *ContractLog) ColBlockEnd() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldBlockEnd())
}

func (*ContractLog) FieldBlockEnd() string {
	return "BlockEnd"
}

func (m *ContractLog) ColTopic0() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldTopic0())
}

func (*ContractLog) FieldTopic0() string {
	return "Topic0"
}

func (m *ContractLog) ColTopic1() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldTopic1())
}

func (*ContractLog) FieldTopic1() string {
	return "Topic1"
}

func (m *ContractLog) ColTopic2() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldTopic2())
}

func (*ContractLog) FieldTopic2() string {
	return "Topic2"
}

func (m *ContractLog) ColTopic3() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldTopic3())
}

func (*ContractLog) FieldTopic3() string {
	return "Topic3"
}

func (m *ContractLog) ColCreatedAt() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldCreatedAt())
}

func (*ContractLog) FieldCreatedAt() string {
	return "CreatedAt"
}

func (m *ContractLog) ColUpdatedAt() *builder.Column {
	return ContractLogTable.ColByFieldName(m.FieldUpdatedAt())
}

func (*ContractLog) FieldUpdatedAt() string {
	return "UpdatedAt"
}

func (m *ContractLog) CondByValue(db sqlx.DBExecutor) builder.SqlCondition {
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

func (m *ContractLog) Create(db sqlx.DBExecutor) error {

	if m.CreatedAt.IsZero() {
		m.CreatedAt.Set(time.Now())
	}

	if m.UpdatedAt.IsZero() {
		m.UpdatedAt.Set(time.Now())
	}

	_, err := db.Exec(sqlx.InsertToDB(db, m, nil))
	return err
}

func (m *ContractLog) List(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) ([]ContractLog, error) {
	var (
		tbl = db.T(m)
		lst = make([]ContractLog, 0)
	)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ContractLog.List")}, adds...)
	err := db.QueryAndScan(builder.Select(nil).From(tbl, adds...), &lst)
	return lst, err
}

func (m *ContractLog) Count(db sqlx.DBExecutor, cond builder.SqlCondition, adds ...builder.Addition) (cnt int64, err error) {
	tbl := db.T(m)
	adds = append([]builder.Addition{builder.Where(cond), builder.Comment("ContractLog.List")}, adds...)
	err = db.QueryAndScan(builder.Select(builder.Count()).From(tbl, adds...), &cnt)
	return
}

func (m *ContractLog) FetchByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ContractLog.FetchByID"),
			),
		m,
	)
	return err
}

func (m *ContractLog) FetchByContractLogID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ContractLogID").Eq(m.ContractLogID),
					),
				),
				builder.Comment("ContractLog.FetchByContractLogID"),
			),
		m,
	)
	return err
}

func (m *ContractLog) FetchByProjectNameAndEventTypeAndChainIDAndContractAddressAndTopic0AndTopic1AndTopic2AndTopic3AndUniq(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	err := db.QueryAndScan(
		builder.Select(nil).
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectName").Eq(m.ProjectName),
						tbl.ColByFieldName("EventType").Eq(m.EventType),
						tbl.ColByFieldName("ChainID").Eq(m.ChainID),
						tbl.ColByFieldName("ContractAddress").Eq(m.ContractAddress),
						tbl.ColByFieldName("Topic0").Eq(m.Topic0),
						tbl.ColByFieldName("Topic1").Eq(m.Topic1),
						tbl.ColByFieldName("Topic2").Eq(m.Topic2),
						tbl.ColByFieldName("Topic3").Eq(m.Topic3),
						tbl.ColByFieldName("Uniq").Eq(m.Uniq),
					),
				),
				builder.Comment("ContractLog.FetchByProjectNameAndEventTypeAndChainIDAndContractAddressAndTopic0AndTopic1AndTopic2AndTopic3AndUniq"),
			),
		m,
	)
	return err
}

func (m *ContractLog) UpdateByIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

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
				builder.Comment("ContractLog.UpdateByIDWithFVs"),
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

func (m *ContractLog) UpdateByID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByIDWithFVs(db, fvs)
}

func (m *ContractLog) UpdateByContractLogIDWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ContractLogID").Eq(m.ContractLogID),
				),
				builder.Comment("ContractLog.UpdateByContractLogIDWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByContractLogID(db)
	}
	return nil
}

func (m *ContractLog) UpdateByContractLogID(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByContractLogIDWithFVs(db, fvs)
}

func (m *ContractLog) UpdateByProjectNameAndEventTypeAndChainIDAndContractAddressAndTopic0AndTopic1AndTopic2AndTopic3AndUniqWithFVs(db sqlx.DBExecutor, fvs builder.FieldValues) error {

	if _, ok := fvs["UpdatedAt"]; !ok {
		fvs["UpdatedAt"] = types.Timestamp{Time: time.Now()}
	}
	tbl := db.T(m)
	res, err := db.Exec(
		builder.Update(tbl).
			Where(
				builder.And(
					tbl.ColByFieldName("ProjectName").Eq(m.ProjectName),
					tbl.ColByFieldName("EventType").Eq(m.EventType),
					tbl.ColByFieldName("ChainID").Eq(m.ChainID),
					tbl.ColByFieldName("ContractAddress").Eq(m.ContractAddress),
					tbl.ColByFieldName("Topic0").Eq(m.Topic0),
					tbl.ColByFieldName("Topic1").Eq(m.Topic1),
					tbl.ColByFieldName("Topic2").Eq(m.Topic2),
					tbl.ColByFieldName("Topic3").Eq(m.Topic3),
					tbl.ColByFieldName("Uniq").Eq(m.Uniq),
				),
				builder.Comment("ContractLog.UpdateByProjectNameAndEventTypeAndChainIDAndContractAddressAndTopic0AndTopic1AndTopic2AndTopic3AndUniqWithFVs"),
			).
			Set(tbl.AssignmentsByFieldValues(fvs)...),
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return m.FetchByProjectNameAndEventTypeAndChainIDAndContractAddressAndTopic0AndTopic1AndTopic2AndTopic3AndUniq(db)
	}
	return nil
}

func (m *ContractLog) UpdateByProjectNameAndEventTypeAndChainIDAndContractAddressAndTopic0AndTopic1AndTopic2AndTopic3AndUniq(db sqlx.DBExecutor, zeros ...string) error {
	fvs := builder.FieldValueFromStructByNoneZero(m, zeros...)
	return m.UpdateByProjectNameAndEventTypeAndChainIDAndContractAddressAndTopic0AndTopic1AndTopic2AndTopic3AndUniqWithFVs(db, fvs)
}

func (m *ContractLog) Delete(db sqlx.DBExecutor) error {
	_, err := db.Exec(
		builder.Delete().
			From(
				db.T(m),
				builder.Where(m.CondByValue(db)),
				builder.Comment("ContractLog.Delete"),
			),
	)
	return err
}

func (m *ContractLog) DeleteByID(db sqlx.DBExecutor) error {
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
				builder.Comment("ContractLog.DeleteByID"),
			),
	)
	return err
}

func (m *ContractLog) DeleteByContractLogID(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ContractLogID").Eq(m.ContractLogID),
					),
				),
				builder.Comment("ContractLog.DeleteByContractLogID"),
			),
	)
	return err
}

func (m *ContractLog) DeleteByProjectNameAndEventTypeAndChainIDAndContractAddressAndTopic0AndTopic1AndTopic2AndTopic3AndUniq(db sqlx.DBExecutor) error {
	tbl := db.T(m)
	_, err := db.Exec(
		builder.Delete().
			From(
				tbl,
				builder.Where(
					builder.And(
						tbl.ColByFieldName("ProjectName").Eq(m.ProjectName),
						tbl.ColByFieldName("EventType").Eq(m.EventType),
						tbl.ColByFieldName("ChainID").Eq(m.ChainID),
						tbl.ColByFieldName("ContractAddress").Eq(m.ContractAddress),
						tbl.ColByFieldName("Topic0").Eq(m.Topic0),
						tbl.ColByFieldName("Topic1").Eq(m.Topic1),
						tbl.ColByFieldName("Topic2").Eq(m.Topic2),
						tbl.ColByFieldName("Topic3").Eq(m.Topic3),
						tbl.ColByFieldName("Uniq").Eq(m.Uniq),
					),
				),
				builder.Comment("ContractLog.DeleteByProjectNameAndEventTypeAndChainIDAndContractAddressAndTopic0AndTopic1AndTopic2AndTopic3AndUniq"),
			),
	)
	return err
}
