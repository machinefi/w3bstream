package example

import (
	"database/sql/driver"

	"github.com/machinefi/w3bstream/pkg/depends/base/types"
)

// @def primary ID
// @def index        I_nickname/BTREE Name
// @def index        I_username       Username
// @def index        I_geom/SPATIAL   (#Geom)
// @def unique_index UI_name          Name
// @def unique_index UI_id_org        ID OrgID

// User 用户表
type User struct {
	ID uint64 `db:"f_id,autoincrement"`

	Name     string     `db:"f_name,default=''"`     // 姓名
	Nickname string     `db:"f_nickname,default=''"` // 昵称
	Username string     `db:"f_username,default=''"` // 用户名
	Gender   Gender     `db:"f_gender,default='0'"`
	Boolean  bool       `db:"f_boolean,default=false"`
	Geom     GeomString `db:"f_geom"`
	// @rel Org.ID
	// 关联组织
	// 组织ID
	OrgID     uint64          `db:"f_org_id"`
	CreatedAt types.Timestamp `db:"f_created_at,default='0'"`
	UpdatedAt types.Timestamp `db:"f_updated_at,default='0'"`
	DeletedAt types.Timestamp `db:"f_deleted_at,default='0'"`
}

type GeomString struct {
	V string
}

func (g GeomString) Value() (driver.Value, error) {
	return g.V, nil
}

func (g *GeomString) Scan(src interface{}) error {
	return nil
}

func (GeomString) DataType(driverName string) string {
	if driverName == "mysql" {
		return "geometry"
	}
	return "geometry(Point)"
}

func (GeomString) ValueEx() string {
	return "ST_GeomFromText(?)"
}
