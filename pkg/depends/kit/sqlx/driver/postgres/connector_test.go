package postgres_test

import (
	"context"
	"testing"

	"github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/driver/postgres"
	"github.com/machinefi/w3bstream/pkg/depends/testutil/buildertestutil"
)

func TestConnector(t *testing.T) {
	c := &postgres.Connector{}

	table := builder.T("t",
		builder.Col("F_id").Type(uint64(0), ",autoincrement"),
		builder.Col("f_old_name").Type("", ",deprecated=f_name"),
		builder.Col("f_name").Type("", ",size=128,default=''"),
		builder.Col("F_geo").Type(&buildertestutil.Point{}, ""),
		builder.Col("F_created_at").Type(int64(0), ",default='0'"),
		builder.Col("F_updated_at").Type(int64(0), ",default='0'"),
		builder.PrimaryKey(builder.Cols("F_id")),
		builder.UniqueIndex("I_name", builder.Cols("F_name")).Using("BTREE"),
		builder.Index("I_created_at", builder.Cols("F_created_at")).Using("BTREE"),
		builder.Index("I_geo", builder.Cols("F_geo")).Using("SPATIAL"),
	)

	cases := map[string]struct {
		expr   builder.SqlExpr
		expect builder.SqlExpr
	}{
		"CreateDatabase": {
			c.CreateDatabase("db"),
			builder.Expr( /* language=PostgreSQL */ `CREATE DATABASE db;`),
		},
		"DropDatabase": {
			c.DropDatabase("db"),
			builder.Expr( /* language=PostgreSQL */ `DROP DATABASE IF EXISTS db;`),
		},
		"AddIndex": {
			c.AddIndex(table.Key("I_name")),
			builder.Expr( /* language=PostgreSQL */ "CREATE UNIQUE INDEX t_i_name ON t USING BTREE (f_name);"),
		},
		"AddPrimaryKey": {
			c.AddIndex(table.Key("PRIMARY")),
			builder.Expr( /* language=PostgreSQL */ "ALTER TABLE t ADD PRIMARY KEY (f_id);"),
		},
		"AddSpatialIndex": {
			c.AddIndex(table.Key("I_geo")),
			builder.Expr( /* language=PostgreSQL */ "CREATE INDEX t_i_geo ON t USING GIST (f_geo);"),
		},
		"DropIndex": {
			c.DropIndex(table.Key("I_name")),
			builder.Expr( /* language=PostgreSQL */ "DROP INDEX IF EXISTS t_i_name"),
		},
		"DropPrimaryKey": {
			c.DropIndex(table.Key("PRIMARY")),
			builder.Expr( /* language=PostgreSQL */ "ALTER TABLE t DROP CONSTRAINT t_pkey;"),
		},
		"CreateTableIsNotExists": {
			c.CreateTableIsNotExists(table)[0],
			builder.Expr( /* language=PostgreSQL */ `CREATE TABLE IF NOT EXISTS t (
	f_id bigserial NOT NULL,
	f_name varchar(128) NOT NULL DEFAULT '',
	f_geo POINT NOT NULL,
	f_created_at bigint NOT NULL DEFAULT '0',
	f_updated_at bigint NOT NULL DEFAULT '0',
	PRIMARY KEY (f_id)
);`),
		},
		"DropTable": {
			c.DropTable(table),
			builder.Expr( /* language=PostgreSQL */ "DROP TABLE IF EXISTS t;"),
		},
		"TruncateTable": {
			c.TruncateTable(table),
			builder.Expr( /* language=PostgreSQL */ "TRUNCATE TABLE t;"),
		},
		"AddColumn": {
			c.AddColumn(table.Col("F_name")),
			builder.Expr( /* language=PostgreSQL */ "ALTER TABLE t ADD COLUMN f_name varchar(128) NOT NULL DEFAULT '';"),
		},
		"DropColumn": {
			c.DropColumn(table.Col("F_name")),
			builder.Expr( /* language=PostgreSQL */ "ALTER TABLE t DROP COLUMN f_name;"),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			gomega.NewWithT(t).Expect(c.expr).
				To(buildertestutil.BeExpr(c.expr.Ex(context.Background()).Query()))
		})
	}
}
