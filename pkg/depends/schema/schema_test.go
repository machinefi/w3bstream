package schema_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/machinefi/w3bstream/pkg/depends/schema"
)

func ExampleSchema() {
	s := schema.Schema{
		Name: "demo",
		Tables: []*schema.Table{
			{
				Name: "tbl",
				Desc: "test table",
				Cols: []*schema.Column{{
					Name: "f_username",
					Constrains: &schema.ColumnType{
						Datatype: schema.DATATYPE__TEXT,
						Length:   255,
						Desc:     "user name",
					},
				}, {
					Name: "f_gender",
					Constrains: &schema.ColumnType{
						Datatype: schema.DATATYPE__UINT8,
						Length:   255,
						Desc:     "user name",
					},
				}},
				Keys: []*schema.Key{{
					Name:     "ui_username",
					IsUnique: true,
					IndexDef: schema.IndexDef{ColumnNames: []string{"f_username"}},
				}},
				WithSoftDeletion: true,
				WithPrimaryKey:   true,
			},
		},
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	data, err = json.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	// Output:
	// {
	//   "name": "demo",
	//   "tables": [
	//     {
	//       "name": "tbl",
	//       "desc": "test table",
	//       "cols": [
	//         {
	//           "name": "f_username",
	//           "constrains": {
	//             "datatype": "TEXT",
	//             "length": 255,
	//             "desc": "user name"
	//           }
	//         },
	//         {
	//           "name": "f_gender",
	//           "constrains": {
	//             "datatype": "UINT8",
	//             "length": 255,
	//             "desc": "user name"
	//           }
	//         }
	//       ],
	//       "keys": [
	//         {
	//           "name": "ui_username",
	//           "isUnique": true,
	//           "columnNames": [
	//             "f_username"
	//           ]
	//         }
	//       ],
	//       "withSoftDeletion": true,
	//       "withPrimaryKey": true
	//     }
	//   ]
	// }
	// {"name":"demo","tables":[{"name":"tbl","desc":"test table","cols":[{"name":"f_username","constrains":{"datatype":"TEXT","length":255,"desc":"user name"}},{"name":"f_gender","constrains":{"datatype":"UINT8","length":255,"desc":"user name"}}],"keys":[{"name":"ui_username","isUnique":true,"columnNames":["f_username"]}],"withSoftDeletion":true,"withPrimaryKey":true}]}
}

var (
	Schema       *schema.Schema
	SchemaConfig = []byte(`{
  "name": "demo",
  "tables": [
    {
      "name": "tbl",
      "desc": "test table",
      "cols": [
        {
          "name": "f_username",
          "constrains": {
            "datatype": "TEXT",
            "length": 255,
            "desc": "user name"
          }
        },
        {
          "name": "f_gender",
          "constrains": {
            "datatype": "UINT8",
            "length": 255,
            "default": "0",
            "desc": "user name"
          }
        }
      ],
      "keys": [
        {
          "name": "ui_username",
          "isUnique": true,
          "columnNames": [
            "f_username"
          ]
        }
      ],
      "withSoftDeletion": true,
      "withPrimaryKey": true
    }
  ]
}`)
)

func init() {
	var err error
	Schema, err = schema.FromConfig(SchemaConfig)
	if err != nil {
		panic(err)
	}
	Schema.WithName("demo")
	if err = Schema.Init(); err != nil {
		panic(err)
	}
}

func ExampleSchema_CreateSchema() {
	q := Schema.CreateSchema()
	fmt.Println(q.Ex(context.Background()).Query())

	// Output:
	// CREATE SCHEMA IF NOT EXISTS wasm_project__demo;
}
