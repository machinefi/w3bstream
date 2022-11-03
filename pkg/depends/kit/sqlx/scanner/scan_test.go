package scanner_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/scanner"
)

type T struct {
	I int    `db:"f_i"`
	S string `db:"f_s"`
}

type Any string

type T2 T

func (t *T2) ColumnReceivers() map[string]interface{} {
	return map[string]interface{}{
		"f_i": &t.I,
		"f_s": &t.S,
	}
}

type TDataList struct {
	Data []T
}

func (*TDataList) New() interface{} { return &T{} }

func (l *TDataList) Next(v interface{}) error {
	t := v.(*T)
	l.Data = append(l.Data, *t)
	return nil
}

func BenchmarkScan(b *testing.B) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	b.Run("ScanToStruct", func(b *testing.B) {
		query := "SELECT f_i,f_s FROM t"
		_ = mock.ExpectQuery(query).
			WillReturnRows(
				mock.NewRows([]string{"f_i", "f_s"}).AddRow(2, "4"),
			)
		target := &T{}
		for i := 0; i < b.N; i++ {
			rows, _ := db.Query(query)
			_ = Scan(context.Background(), rows, target)
		}
		b.Log(target)
	})
	b.Run("ScanToColumnReceivers", func(b *testing.B) {
		sql := "SELECT f_i,f_s from t"

		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "4")

		_ = mock.ExpectQuery(sql).WillReturnRows(mockRows)

		target := &T2{}

		for i := 0; i < b.N; i++ {
			rows, _ := db.Query(sql)
			_ = Scan(context.Background(), rows, target)
		}

		b.Log(target)
	})
}

func TestScan(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	t.Run("ScanToStruct", func(t *testing.T) {
		sql := "SELECT f_i,f_s from t"

		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "4")

		_ = mock.ExpectQuery(sql).WillReturnRows(mockRows)

		target := &T{}
		rows, _ := db.Query(sql)
		err := Scan(context.Background(), rows, target)
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
		gomega.NewWithT(t).Expect(target).To(gomega.Equal(&T{
			I: 2,
			S: "4",
		}))
	})

	t.Run("ScanToStructWithColumnReceivers", func(t *testing.T) {
		sql := "SELECT f_i,f_s from t"

		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "4")

		_ = mock.ExpectQuery(sql).WillReturnRows(mockRows)

		target := &T2{}
		rows, _ := db.Query(sql)
		err := Scan(context.Background(), rows, target)
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
		gomega.NewWithT(t).Expect(target).To(gomega.Equal(&T2{
			I: 2,
			S: "4",
		}))
	})

	t.Run("ScanToStructWithoutNoRecord", func(t *testing.T) {
		sql := "SELECT f_i,f_s from t"

		mockRows := mock.NewRows([]string{"f_i", "f_s"})

		_ = mock.ExpectQuery(sql).WillReturnRows(mockRows)

		target := &T{}
		rows, err := db.Query(sql)
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

		err = Scan(context.Background(), rows, target)
		gomega.NewWithT(t).Expect(err).To(gomega.Equal(RecordNotFound))
	})

	t.Run("Scan to count", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"count(1)"})
		mockRows.AddRow(10)

		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		count := 0
		rows, err := db.Query("SELECT count(1) from t")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

		err = Scan(context.Background(), rows, &count)
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
		gomega.NewWithT(t).Expect(count).To(gomega.Equal(10))
	})

	t.Run("ScanToCountFailedWhenBadReceiver", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"count(1)"})
		mockRows.AddRow(10)

		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		v := Any("")
		rows, err := db.Query("SELECT count(1) from t")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

		err = Scan(context.Background(), rows, &v)
		gomega.NewWithT(t).Expect(err).NotTo(gomega.BeNil())
	})

	t.Run("ScanToSlice", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "2")
		mockRows.AddRow(3, "3")

		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		list := make([]T, 0)
		rows, err := db.Query("SELECT f_i,f_b from t")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

		err = Scan(context.Background(), rows, &list)

		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
		gomega.NewWithT(t).Expect(list).To(gomega.Equal([]T{
			{
				I: 2,
				S: "2",
			},
			{
				I: 3,
				S: "3",
			},
		}))
	})

	t.Run("ScanToIter", func(t *testing.T) {
		mockRows := mock.NewRows([]string{"f_i", "f_s"})
		mockRows.AddRow(2, "2")
		mockRows.AddRow(3, "3")

		_ = mock.ExpectQuery("SELECT .+ from t").WillReturnRows(mockRows)

		rows, err := db.Query("SELECT f_i,f_b from t")
		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

		list := TDataList{}

		err = Scan(context.Background(), rows, &list)

		gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
		gomega.NewWithT(t).Expect(list.Data).To(gomega.Equal([]T{
			{
				I: 2,
				S: "2",
			},
			{
				I: 3,
				S: "3",
			},
		}))
	})
}
