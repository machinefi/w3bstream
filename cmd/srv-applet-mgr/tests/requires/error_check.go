package requires

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
)

func CheckError(t *testing.T, err error, expect error) {
	if err == nil && expect == nil {
		return
	}
	se1 := statusx.FromErr(err)
	se2 := statusx.FromErr(expect)

	NewWithT(t).Expect(se1.Key).To(Equal(se2.Key))
}
