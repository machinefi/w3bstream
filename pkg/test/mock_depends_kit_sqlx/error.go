package mock_sqlx

import (
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
)

var (
	ErrConflict = sqlx.NewSqlError(sqlx.SqlErrTypeConflict, "")
	ErrNotFound = sqlx.NewSqlError(sqlx.SqlErrTypeNotFound, "")
	ErrDatabase = errors.New("database error")
)

func ExpectError(t *testing.T, err error, se status.Error) {
	NewWithT(t).Expect(err).NotTo(BeNil())
	expect, ok := statusx.IsStatusErr(err)
	NewWithT(t).Expect(ok).To(BeTrue())
	NewWithT(t).Expect(expect.Key).To(Equal(se.Key()))
}
