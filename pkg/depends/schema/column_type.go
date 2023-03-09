package schema

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
)

// ColumnType defines universal column constrains
type ColumnType struct {
	Datatype      Datatype `json:"datatype"`
	Length        uint64   `json:"length,omitempty"`
	Decimal       uint64   `json:"decimal,omitempty"`
	Default       *string  `json:"default,omitempty"`
	Null          bool     `json:"null,omitempty"`
	AutoIncrement bool     `json:"autoincrement,omitempty"`
	Desc          string   `json:"desc,omitempty"`
}

var _ builder.SqlExpr = (*ColumnType)(nil)

func (t *ColumnType) IsNil() bool {
	return t == nil || t.Datatype == DATATYPE_UNKNOWN
}

func (t *ColumnType) Ex(ctx context.Context) *builder.Ex {
	e := builder.Expr("")
	e.WriteQuery(t.AutoCompleteLengthDatatype())
	e.WriteQuery(t.TypeModify())
	return e
}

func (t *ColumnType) DatabaseDatatype() string {
	datatype := ""
	switch dt := t.Datatype; dt {
	case
		DATATYPE__INT, DATATYPE__INT8, DATATYPE__INT16, DATATYPE__INT32,
		DATATYPE__UINT, DATATYPE__UINT8, DATATYPE__UINT16, DATATYPE__UINT32:
		if t.AutoIncrement {
			datatype = "serial"
		} else {
			datatype = "integer"
		}
	case DATATYPE__INT64, DATATYPE__UINT64:
		if t.AutoIncrement {
			datatype = "bigserial"
		} else {
			datatype = "bigint"
		}
	case DATATYPE__FLOAT32:
		datatype = "real"
	case DATATYPE__FLOAT64:
		datatype = "double precision"
	case DATATYPE__TEXT:
		if t.Length < 65536/3 {
			datatype = "character varying"
		} else {
			datatype = "text"
		}
	case DATATYPE__BOOL:
		datatype = "boolean"
	case DATATYPE__TIMESTAMP:
		// TODO should use "timestamp without time zone"
		datatype = "bigint"
	case DATATYPE__DECIMAL:
		datatype = "decimal"
	case DATATYPE__NUMERIC:
		datatype = "numeric"
	default:
		panic(fmt.Errorf("unsupport type: %v", dt))
	}
	return datatype
}

func (t *ColumnType) AutoCompleteLengthDatatype() string {
	dt := t.DatabaseDatatype()
	if dt == "character varying" || dt == "real" ||
		dt == "double precision" || dt == "decimal" || dt == "numeric" {
		size := t.Length
		if dt == "character varying" {
			if size == 0 {
				size = 255
			}
		}
		if size > 0 {
			sizestr := strconv.FormatUint(size, 10)
			if t.Decimal > 0 {
				dt += "(" + sizestr + "," + strconv.FormatUint(t.Decimal, 10) + ")"
			} else {
				dt += "(" + sizestr + ")"
			}
		}
	}
	return dt
}

func (t *ColumnType) TypeModify() string {
	b := bytes.NewBuffer(nil)

	if !t.Null {
		b.WriteString(" NOT NULL")
	}
	if t.Default != nil && len(*t.Default) > 0 {
		b.WriteString(" DEFAULT ")
		b.WriteString(t.NormalizeDefaultValue())
	}
	return b.String()
}

func (t *ColumnType) NormalizeDefaultValue() string {
	if t.Default == nil || len(*t.Default) == 0 {
		return ""
	}
	dv := *t.Default
	if dv[0] == '\'' {
		if strings.Contains(dv, "'::") {
			return dv
		}
		return dv + "::" + t.DatabaseDatatype()
	}
	return dv
}
