package schema

import "github.com/machinefi/w3bstream/pkg/depends/gen/codegen"

//go:generate toolkit gen enum Datatype
type Datatype uint8

const (
	DATATYPE_UNKNOWN Datatype = iota
	DATATYPE__INT
	DATATYPE__INT8
	DATATYPE__INT16
	DATATYPE__INT32
	DATATYPE__INT64
	DATATYPE__UINT
	DATATYPE__UINT8
	DATATYPE__UINT16
	DATATYPE__UINT32
	DATATYPE__UINT64
	DATATYPE__STRING
	DATATYPE__BOOL
	DATATYPE__TIMESTAMP
)

func (d *Datatype) CodeGenType(f *codegen.File) codegen.SnippetType {
	switch *d {
	case DATATYPE__INT:
		return codegen.Int
	case DATATYPE__INT8:
		return codegen.Int8
	case DATATYPE__INT16:
		return codegen.Int16
	case DATATYPE__INT32:
		return codegen.Int32
	case DATATYPE__INT64:
		return codegen.Int64
	case DATATYPE__UINT:
		return codegen.Uint
	case DATATYPE__UINT8:
		return codegen.Uint8
	case DATATYPE__UINT16:
		return codegen.Uint16
	case DATATYPE__UINT32:
		return codegen.Uint32
	case DATATYPE__UINT64:
		return codegen.Uint64
	case DATATYPE__STRING:
		return codegen.String
	case DATATYPE__BOOL:
		return codegen.Bool
	case DATATYPE__TIMESTAMP:
		return codegen.Type(f.Use(PkgTypes, "Timestamp"))
	}
	return codegen.Int
}

var PkgTypes = "github.com/machinefi/w3bstream/pkg/depends/base/types"
