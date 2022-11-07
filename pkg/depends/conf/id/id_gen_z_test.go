package id_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/pkg/depends/conf/id"
)

func TestFromIP(t *testing.T) {
	g, err := id.FromLocalIP()
	NewWithT(t).Expect(err).To(BeNil())

	sfid, err := g.ID()
	NewWithT(t).Expect(err).To(BeNil())
	t.Log(sfid)
}
