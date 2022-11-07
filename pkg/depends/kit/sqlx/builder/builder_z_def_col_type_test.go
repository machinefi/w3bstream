package builder_test

import (
	"reflect"
	"testing"

	g "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/depends/x/typesx"
)

func TestAnalyzeColumnType(t *testing.T) {
	cases := []struct {
		name string
		tag  string
		val  *ColumnType
	}{
		{
			"Deprecated",
			`,deprecated=deprecated_f_old_name`,
			&ColumnType{
				Type:           typesx.FromReflectType(reflect.TypeOf(1)),
				DeprecatedActs: &DeprecatedActs{RenameTo: "deprecated_f_old_name"},
			},
		},
		{
			"AutoIncrement",
			`,autoincrement`,
			&ColumnType{
				Type:          typesx.FromReflectType(reflect.TypeOf(1)),
				AutoIncrement: true,
			},
		}, {

			"Null",
			`,null`,
			&ColumnType{
				Type: typesx.FromReflectType(reflect.TypeOf(float64(0))),
				Null: true,
			},
		}, {
			"Size",
			`,size=2`,
			&ColumnType{
				Type:   typesx.FromReflectType(reflect.TypeOf("")),
				Length: 2,
			},
		}, {
			"Decimal",
			`,decimal=1`,
			&ColumnType{
				Type:    typesx.FromReflectType(reflect.TypeOf(float64(0))),
				Decimal: 1,
			},
		}, {
			"Default",
			`,default='1'`,
			&ColumnType{
				Type:    typesx.FromReflectType(reflect.TypeOf("")),
				Default: ptrx.String(`'1'`),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g.NewWithT(t).
				Expect(AnalyzeColumnType(c.val.Type, c.tag)).
				To(g.Equal(c.val))
		})
	}
}
