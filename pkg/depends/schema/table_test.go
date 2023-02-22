package schema_test

import (
	"context"
	"fmt"
)

func ExampleTable_CreateIfNotExists() {
	t := Schema.T("tbl")
	es := t.CreateIfNotExists()
	for _, e := range es {
		fmt.Println(e.Ex(context.Background()).Query())
	}

	// Output:
	// CREATE TABLE IF NOT EXISTS wasm_project__demo.tbl (
	// 	f_id bigserial NOT NULL,
	// 	f_username character varying(255) NOT NULL,
	// 	f_gender integer NOT NULL DEFAULT '0'::integer,
	// 	f_created_at bigint NOT NULL DEFAULT '0'::bigint,
	// 	f_updated_at bigint NOT NULL DEFAULT '0'::bigint,
	// 	f_deleted_at bigint NOT NULL DEFAULT '0'::bigint
	// );
	// CREATE UNIQUE INDEX IF NOT EXISTS tbl_ui_username ON wasm_project__demo.tbl (f_username,f_deleted_at);
}
