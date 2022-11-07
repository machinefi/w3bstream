package builder_test

import (
	"context"
	"testing"

	g "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	. "github.com/machinefi/w3bstream/pkg/depends/testutil/buildertestutil"
)

func TestAssignment(t *testing.T) {
	t.Run("ColumnsAndValues", func(t *testing.T) {
		g.NewWithT(t).Expect(
			ColumnsAndValues(Cols("a", "b"), 1, 2, 3, 4).
				Ex(ContextWithToggleUseValues(context.Background(), true)),
		).To(BeExpr("(a,b) VALUES (?,?),(?,?)", 1, 2, 3, 4))
	})
}
