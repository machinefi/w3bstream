package wasm_test

/*
import (
	"context"
	"testing"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/global"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
	. "github.com/onsi/gomega"
)

func TestDatabase_Init(t *testing.T) {
	ctx := global.WithContext(context.Background())
	ctx = types.WithProject(ctx,
		&models.Project{RelProject: models.RelProject{ProjectID: 1234567}},
	)
	// if only want to inspect queries, effect code next line
	// ctx = migration.WithInspectionOutput(ctx, os.Stderr)
	database := &wasm.Database{
		Schemas: []*wasm.Schema{
			{
				Tables: []*wasm.Table{
					{
						Name: "t_demo",
						Desc: "demo table",
						Cols: []*wasm.Column{
							{
								Name: "f_id",
								Constrains: wasm.Constrains{
									Datatype:      enums.WASM_DB_DATATYPE__INT64,
									AutoIncrement: true,
									Desc:          "primary id",
								},
							},
							{
								Name: "f_name",
								Constrains: wasm.Constrains{
									Datatype: enums.WASM_DB_DATATYPE__TEXT,
									Length:   255,
									Desc:     "name",
								},
							},
							{
								Name: "f_amount",
								Constrains: wasm.Constrains{
									Datatype: enums.WASM_DB_DATATYPE__FLOAT64,
									Desc:     "amount",
								},
							},
							{
								Name: "f_income",
								Constrains: wasm.Constrains{
									Datatype: enums.WASM_DB_DATATYPE__DECIMAL,
									Default:  ptrx.Ptr("0"),
									Length:   128,
									Decimal:  512,
									Desc:     "income",
								},
							},
							{
								Name: "f_comment",
								Constrains: wasm.Constrains{
									Datatype: enums.WASM_DB_DATATYPE__TEXT,
									Default:  ptrx.Ptr(""),
									Null:     true,
									Desc:     "comment",
								},
							},
						},
						Keys: nil,
					},
				},
			},
		},
	}

	// migration test
	err := database.Init(ctx)
	NewWithT(t).Expect(err).To(BeNil())

	d, err := database.WithDefaultSchema()
	NewWithT(t).Expect(err).To(BeNil())

	// fetch column value
	_, err = d.Exec(builder.Expr("SELECT f_id FROM t_demo"))
	NewWithT(t).Expect(err).To(BeNil())
	// fetch a nonexistence column
	_, err = d.Exec(builder.Expr("SELECT f_xxx FROM t_demo"))
	NewWithT(t).Expect(err).NotTo(BeNil())

	// migration test: add table
	table := *(database.Schemas[0].Tables[0])
	table.Name = "f_demo_2"
	database.Schemas[0].Tables = append(database.Schemas[0].Tables, &table)

	err = database.Init(ctx)
	NewWithT(t).Expect(err).To(BeNil())
	_, err = d.Exec(builder.Expr("SELECT f_id FROM f_demo_2"))
	NewWithT(t).Expect(err).To(BeNil())
	// fetch a nonexistence column
	_, err = d.Exec(builder.Expr("SELECT f_xxx FROM f_demo_2"))
	NewWithT(t).Expect(err).NotTo(BeNil())

	// migration test: add column
	table2 := database.Schemas[0].Tables[1]
	table2.Cols = append(table2.Cols, &wasm.Column{
		Name: "f_added_after",
		Constrains: wasm.Constrains{
			Datatype: enums.WASM_DB_DATATYPE__TEXT,
			Default:  ptrx.Ptr(""),
			Null:     true,
			Desc:     "comment",
		},
	})
	err = database.Init(ctx)
	NewWithT(t).Expect(err).To(BeNil())
	// fetch column f_added_after
	_, err = d.Exec(builder.Expr("SELECT f_added_after FROM f_demo_2"))
	NewWithT(t).Expect(err).To(BeNil())

	// migration test: drop column
	// SHOULD NOT support drop column

}
*/
