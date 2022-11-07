package builder_test

import (
	"testing"

	g "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	. "github.com/machinefi/w3bstream/pkg/depends/testutil/buildertestutil"
)

func BenchmarkCols(b *testing.B) {
	columns := Columns{}

	columns.Add(
		Col("f_id").Field("ID").Type(1, `,autoincrement`),
		Col("f_name").Field("Name").Type(1, ``),
		Col("f_f1").Field("F1").Type(1, ``),
		Col("f_f2").Field("F2").Type(1, ``),
		Col("f_f3").Field("F3").Type(1, ``),
		Col("f_f4").Field("F4").Type(1, ``),
		Col("f_f5").Field("F5").Type(1, ``),
		Col("f_f6").Field("F6").Type(1, ``),
		Col("f_f7").Field("F7").Type(1, ``),
		Col("f_f8").Field("F8").Type(1, ``),
		Col("f_f9").Field("F9").Type(1, ``),
	)

	b.Run("Pick", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = columns.ColByFieldName("F3")
		}
	})
	b.Run("MultiPick", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = columns.ColsByFieldNames("ID", "Name")
		}
	})
	b.Run("PickAll", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = columns.ColsByFieldNames()
		}
	})
}

func TestColumns(t *testing.T) {
	columns := Columns{}

	t.Run("EmptyColumns", func(t *testing.T) {
		g.NewWithT(t).Expect(columns.Len()).To(g.Equal(0))
		g.NewWithT(t).Expect(columns.AutoIncrement()).To(g.BeNil())
	})

	t.Run("AddColumns", func(t *testing.T) {
		columns.Add(
			Col("F_id").Field("ID").Type(1, `,autoincrement`),
		)

		autoIncrementCol := columns.AutoIncrement()
		g.NewWithT(t).Expect(autoIncrementCol).NotTo(g.BeNil())
		g.NewWithT(t).Expect(autoIncrementCol.Name).To(g.Equal("f_id"))

		t.Run("GetByFieldName", func(t *testing.T) {
			g.NewWithT(t).Expect(columns.ColByFieldName("ID2")).To(g.BeNil())
			g.NewWithT(t).Expect(MustCols(columns.ColsByFieldNames("ID2")).Len()).To(g.Equal(0))
			g.NewWithT(t).Expect(MustCols(columns.ColsByFieldNames()).Len()).To(g.Equal(1))
			g.NewWithT(t).Expect(MustCols(columns.ColsByFieldNames("ID2")).List()).To(g.HaveLen(0))
			g.NewWithT(t).Expect(MustCols(columns.ColsByFieldNames()).Len()).To(g.Equal(1))
		})
		t.Run("GetByColName", func(t *testing.T) {
			g.NewWithT(t).Expect(MustCols(columns.Cols("F_id")).Len()).To(g.Equal(1))
			g.NewWithT(t).Expect(MustCols(columns.Cols()).Len()).To(g.Equal(1))
			g.NewWithT(t).Expect(MustCols(columns.Cols()).List()).To(g.HaveLen(1))
			g.NewWithT(t).Expect(MustCols(columns.Cols()).FieldNames()).To(g.Equal([]string{"ID"}))
		})
	})
}
