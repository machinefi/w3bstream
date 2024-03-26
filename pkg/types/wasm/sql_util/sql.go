package sql_util

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

func ParseQuery(data []byte) (prestate string, params []interface{}, err error) {
	if !gjson.ValidBytes(data) {
		err = errors.New("query is invalid")
		return
	}

	res := gjson.ParseBytes(data)
	prestateRes := res.Get("statement")
	paramsRes := res.Get("params")
	if !prestateRes.Exists() || !paramsRes.Exists() {
		err = errors.New("query is invalid")
		return
	}
	prestate = prestateRes.String()

	params = make([]interface{}, 0)
	for _, para := range paramsRes.Array() {
		var res interface{}
		res, err = DecodeQueryParam(&para)
		if err != nil {
			return
		}
		params = append(params, res)
	}

	return
}

func DecodeQueryParam(in *gjson.Result) (ret interface{}, err error) {
	switch {
	case in.Get("int32").Exists():
		ret = int32(in.Get("int32").Int())
	case in.Get("int64").Exists():
		ret = int64(in.Get("int64").Int())
	case in.Get("float32").Exists():
		ret = float32(in.Get("float32").Float())
	case in.Get("float64").Exists():
		ret = float64(in.Get("float64").Float())
	case in.Get("string").Exists():
		ret = in.Get("string").String()
	case in.Get("bool").Exists():
		ret = in.Get("bool").Bool()
	case in.Get("bytes").Exists():
		ret, err = base64.StdEncoding.DecodeString(in.Get("bytes").String())
	case in.Get("time").Exists():
		ret, err = time.Parse(time.RFC3339, in.Get("time").String())
	default:
		err = errors.New("fail to decode the param")
	}
	return
}

func JsonifyRows(rawRows *sql.Rows) ([]byte, error) {
	if rawRows == nil {
		return nil, errors.New("rows are empty")
	}

	columnTypes, err := rawRows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	rows := make([]interface{}, 0)
	for rawRows.Next() {
		scanArgs := make([]interface{}, len(columnTypes))
		for i := range columnTypes {
			switch columnTypes[i].DatabaseTypeName() {
			case "VARCHAR", "TEXT", "CHAR":
				scanArgs[i] = new(sql.NullString)
			case "TIMESTAMP", "TIME", "DATE":
				scanArgs[i] = new(sql.NullTime)
			case "BOOL", "BOOLEAN":
				scanArgs[i] = new(sql.NullBool)
			case "INT", "INTEGER", "SMALLINT", "BIGINT", "INT2", "INT4", "INT8":
				scanArgs[i] = new(sql.NullInt64)
			case "FLOAT", "FLOAT4", "FLOAT8", "DOUBLE", "NUMERIC":
				scanArgs[i] = new(sql.NullFloat64)
			default:
				// fmt.Println(columnTypes[i].DatabaseTypeName(), columnTypes[i].ScanType().Name())
				scanArgs[i] = new(sql.NullString)
			}
		}

		if err := rawRows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		entryMap := make(map[string]interface{}, len(columnTypes))
		for i := 0; i < len(columnTypes); i++ {
			colName := columnTypes[i].Name()
			switch v := scanArgs[i].(type) {
			case *sql.NullBool:
				if !v.Valid {
					entryMap[colName] = nil
					continue
				}
				entryMap[colName], err = v.Value()
			case *sql.NullString:
				if !v.Valid {
					entryMap[colName] = nil
					continue
				}
				entryMap[colName], err = v.Value()
			case *sql.NullFloat64:
				if !v.Valid {
					entryMap[colName] = nil
					continue
				}
				entryMap[colName], err = v.Value()
			case *sql.NullInt64:
				if !v.Valid {
					entryMap[colName] = nil
					continue
				}
				entryMap[colName], err = v.Value()
			// TODO: support time encodings
			// case *sql.NullTime:
			// 	if !v.Valid {
			// 		entryMap[colName] = nil
			// 		continue
			// 	}
			// 	entryMap[colName], err = v.Value()
			default:
				entryMap[colName] = scanArgs[i]
			}
			if err != nil {
				return nil, err
			}
		}
		rows = append(rows, entryMap)
	}
	if len(rows) == 0 {
		return []byte{}, nil
	}
	if len(rows) == 1 {
		return json.Marshal(rows[0])
	}
	return json.Marshal(rows)
}
