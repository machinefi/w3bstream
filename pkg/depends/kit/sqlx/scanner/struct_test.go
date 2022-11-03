package scanner_test

import (
	"context"
	"testing"
	_ "unsafe"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/scanner"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/scanner/nullable"
)

func BenchmarkScanStruct(b *testing.B) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mockRows := mock.NewRows([]string{"f_i", "f_s"})

	b.Run("Direct", func(b *testing.B) {
		target := &T2{}
		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		rows, _ := db.Query("SELECT f_i,f_s from t")

		for i := 0; i < b.N; i++ {
			mockRows.AddRow(i, "a")

			cs, _ := rows.Columns()

			if rows.Next() {
				dest := make([]interface{}, len(cs))
				crs := target.ColumnReceivers()
				for i := range cs {
					dest[i] = nullable.NewNullIgnoreScanner(crs[cs[i]])
				}
				_ = rows.Scan(dest...)
			}
		}

		b.Log(target)
	})

	b.Run("ScanByReflect", func(b *testing.B) {
		target := &T{}
		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		rows, _ := db.Query("SELECT f_i,f_s from t")

		for i := 0; i < b.N; i++ {
			mockRows.AddRow(i, "b")

			if rows.Next() {
				_ = scanner.ScanRows(context.Background(), rows, target)
			}
		}

		b.Log(target)
	})

	b.Run("ScanByColumnReceivers", func(b *testing.B) {
		target := &T2{}
		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		rows, _ := db.Query("SELECT f_i,f_s from t")

		for i := 0; i < b.N; i++ {
			mockRows.AddRow(i, "c")

			if rows.Next() {
				_ = scanner.ScanRows(context.Background(), rows, target)
			}
		}

		b.Log(target)
	})
}
