package tests

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/__test__/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/__test__/requires"
	"github.com/machinefi/w3bstream/pkg/errors/status"
)

func TestProjectAPIs(t *testing.T) {
	// defer requires.Serve()()
	c := requires.AuthClient()
	t.Run("Project", func(t *testing.T) {
		t.Run("#CreateProject", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {
				req := &applet_mgr.CreateProject{}
				req.CreateReq.Name = "demo"
				_, _, err := c.CreateProject(req)
				NewWithT(t).Expect(err).NotTo(BeNil())
				requires.CheckError(t, err, status.InvalidAuthAccountID)
			})
		})
	})
}
