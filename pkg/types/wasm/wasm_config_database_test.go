package wasm_test

import (
	"context"
	"os"
	"testing"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/global"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/migration"
	"github.com/machinefi/w3bstream/pkg/depends/x/ptrx"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
	"github.com/machinefi/w3bstream/pkg/types/wasm"
)

func TestDatabase_Init(t *testing.T) {
	ctx := global.WithContext(context.Background())
	ctx = types.WithProject(ctx,
		&models.Project{RelProject: models.RelProject{ProjectID: 1234567}},
	)
	ctx = migration.WithInspectionOutput(ctx, os.Stderr)
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
								Constrains: &wasm.Constrains{
									Datatype:      enums.WASM_DB_DATATYPE__INT64,
									AutoIncrement: true,
									Desc:          "primary id",
								},
							},
							{
								Name: "f_name",
								Constrains: &wasm.Constrains{
									Datatype: enums.WASM_DB_DATATYPE__TEXT,
									Length:   255,
									Desc:     "name",
								},
							},
							{
								Name: "f_amount",
								Constrains: &wasm.Constrains{
									Datatype: enums.WASM_DB_DATATYPE__FLOAT64,
									Desc:     "amount",
								},
							},
							{
								Name: "f_income",
								Constrains: &wasm.Constrains{
									Datatype: enums.WASM_DB_DATATYPE__DECIMAL,
									Desc:     "amount",
								},
							},
							{
								Name: "f_comment",
								Constrains: &wasm.Constrains{
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

	err := database.Init(ctx)
	t.Log(err)
}
