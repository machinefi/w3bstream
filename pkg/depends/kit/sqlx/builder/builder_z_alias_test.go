package builder_test

import (
	"testing"

	g "github.com/onsi/gomega"

	. "github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	. "github.com/machinefi/w3bstream/pkg/depends/testutil/buildertestutil"
)

func TestAlias(t *testing.T) {
	t.Run("Alias", func(t *testing.T) {
		g.NewWithT(t).Expect(Alias(Expr("f_id"), "id")).
			To(BeExpr("f_id AS id"))
	})
}
