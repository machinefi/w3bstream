package builder

import (
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

type ColumnType struct {
	Type           typesx.Type
	DataType       string
	Length         uint64
	Decimal        uint64
	Default        *string
	OnUpdate       *string
	Null           bool
	AutoIncrement  bool
	Comment        string
	Desc           []string
	Rel            []string
	DeprecatedActs *DeprecatedActs
}

func AnalyzeColumnType(t typesx.Type, tag string) *ColumnType {
	ct := &ColumnType{Type: typesx.DeRef(t)}

	if !strings.Contains(tag, ",") {
		return ct
	}

	for _, flag := range strings.Split(tag, ",") {
		kv := strings.Split(flag, "=")
		switch strings.ToLower(kv[0]) {
		case "null":
			ct.Null = true
		case "autoincrement":
			ct.AutoIncrement = true
		case "size":
			if len(kv) == 1 {
				panic("missing size value")
			}
			length, err := strconv.ParseUint(kv[1], 10, 64)
			if err != nil {
				panic("invalid size value: " + kv[1])
			}
			ct.Length = length
		case "decimal":
			if len(kv) == 1 {
				panic("missing size value")
			}
			decimal, err := strconv.ParseUint(kv[1], 10, 64)
			if err != nil {
				panic("invalid decimal value: " + kv[1])
			}
			ct.Decimal = decimal
		case "default":
			if len(kv) == 1 {
				panic("missing default value")
			}
			ct.Default = &kv[1]
		case "onupdate":
			if len(kv) == 1 {
				panic("missing onupdate value")
			}
			ct.OnUpdate = &kv[1]
		case "deprecated":
			rename := ""
			if len(kv) > 1 {
				rename = kv[1]
			}
			ct.DeprecatedActs = &DeprecatedActs{RenameTo: rename}
		}
	}

	return ct
}

type DeprecatedActs struct {
	RenameTo string `name:"rename"`
	// TODO drop column or other action
}
