package scanner

import (
	"context"
	"database/sql"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/scanner/nullable"
	"github.com/machinefi/w3bstream/pkg/depends/x/reflectx"
)

type RowScanner interface {
	Scan(...interface{}) error
}

type WithColumnReceivers interface {
	ColumnReceivers() map[string]interface{}
}

func ScanRows(ctx context.Context, rows *sql.Rows, v interface{}) error {
	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Ptr {
		return errors.Errorf("scan target must be a ptr")
	}

	if scanner, ok := v.(sql.Scanner); ok {
		return rows.Scan(scanner)
	}

	t = reflectx.DeRef(t)

	switch t.Kind() {
	default:
		return rows.Scan(nullable.NewNullIgnoreScanner(v))
	case reflect.Struct:
		cols, err := rows.Columns()
		if err != nil {
			return err
		}
		colc := len(cols)
		if colc < 1 {
			return nil
		}
		dst := make([]interface{}, colc)
		holder := EmptyScanner(0)

		if with, ok := v.(WithColumnReceivers); ok {
			receivers := with.ColumnReceivers()
			for i, name := range cols {
				if c, ok := receivers[strings.ToLower(name)]; ok {
					dst[i] = nullable.NewNullIgnoreScanner(c)
				} else {
					dst[i] = &holder
				}
			}
			return rows.Scan(dst...)
		}

		indexes := map[string]int{}
		for i, name := range cols {
			indexes[strings.ToLower(name)] = i
			dst[i] = &holder
		}

		builder.ForEachFieldValue(ctx, v, func(fv *builder.FieldValue) {
			if fv.TableName != "" {
				name := fv.TableName + "__" + fv.Field.Name
				if i, ok := indexes[name]; ok && i > -1 {
					dst[i] = nullable.NewNullIgnoreScanner(fv.Value.Addr().Interface())
				}
			}
			if i, ok := indexes[fv.Field.Name]; ok && i > -1 {
				dst[i] = nullable.NewNullIgnoreScanner(fv.Value.Addr().Interface())
			}
		})

		return rows.Scan(dst...)
	}
}

type EmptyScanner int

func (s *EmptyScanner) Scan(_ interface{}) error { return nil }
