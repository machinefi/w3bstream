package util_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/util"
)

var (
	gPassword  = "12345"
	gAccountID = "100"
)

func TestHashOfAccountPassword(t *testing.T) {
	cryptPasswd := util.HashOfAccountPassword(gAccountID, gPassword)
	t.Log(cryptPasswd)
	password, err := util.ExtractRawPasswordByAccountAndPassword(gAccountID, cryptPasswd)
	NewWithT(t).Expect(err).To(BeNil())
	// TODO: fix unit test
	NewWithT(t).Expect(password).To(Equal(""))
}
