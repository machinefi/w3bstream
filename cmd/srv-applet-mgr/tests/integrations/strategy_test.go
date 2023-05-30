package integrations

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/requires"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer"
	"github.com/machinefi/w3bstream/pkg/modules/applet"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestStrategyAPIs(t *testing.T) {
	var (
		client      = requires.AuthClient()
		projectName = "test_strategy_project"
		eventType   = "TYPE0"
		handler     = "start"

		appletID   types.SFID
		strategyID types.SFID
	)

	t.Logf("random a project name: %s, use this name create a project and an applet.", projectName)

	t.Run("PrepareProject", func(t *testing.T) {
		t.Run("#CreateProject", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// create project without user defined config(database/env)
				{
					req := &applet_mgr.CreateProject{}
					req.CreateReq.Name = projectName

					rsp, _, err := client.CreateProject(req)

					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.Name).To(Equal(projectName))
				}
			})
		})
		t.Run("#CreateApplet", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				{
					cwd, err := os.Getwd()
					NewWithT(t).Expect(err).To(BeNil())

					filename := path.Join(cwd, "../testdata/log.wasm")
					appletName := "testApplet"
					wasmName := "test.log"

					req := &applet_mgr.CreateApplet{ProjectName: projectName}
					req.CreateReq.File = transformer.MustNewFileHeader("file", filename, bytes.NewBuffer(code))
					req.CreateReq.Info = applet.Info{
						AppletName: appletName,
						WasmName:   wasmName,
					}

					rsp, _, err := client.CreateApplet(req)

					NewWithT(t).Expect(err).To(BeNil())
					appletID = rsp.AppletID
				}
			})
		})
	})

	t.Logf("random a strategy with EventType and Handler: %s - %s, then create it .", eventType, handler)

	t.Run("Strategy", func(t *testing.T) {
		t.Run("#CreateStrategy", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// create strategy
				{
					req := &applet_mgr.CreateStrategy{ProjectName: projectName}
					req.CreateReq.AppletID = appletID
					req.CreateReq.EventType = eventType
					req.CreateReq.Handler = handler

					rsp, _, err := client.CreateStrategy(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.EventType).To(Equal(eventType))
					strategyID = rsp.StrategyID
				}

				// get strategy
				{
					req := &applet_mgr.GetStrategy{StrategyID: strategyID}
					rsp, _, err := client.GetStrategy(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.EventType).To(Equal(eventType))
					NewWithT(t).Expect(rsp.Handler).To(Equal(handler))
				}

				// update strategy
				{
					updateType := "updatetype"
					updateHandle := "updatehandle"
					req := &applet_mgr.UpdateStrategy{StrategyID: strategyID}
					req.UpdateReq.AppletID = appletID
					req.UpdateReq.EventType = updateType
					req.UpdateReq.Handler = updateHandle
					_, err := client.UpdateStrategy(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// remove strategy
				{
					req := &applet_mgr.RemoveStrategy{StrategyID: strategyID}
					_, err := client.RemoveStrategy(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
		})
	})

	t.Run("BatchStrategies", func(t *testing.T) {
		t.Run("#CreateStrategies", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// prepare data
				num := 5
				{
					for i := 0; i < num; i++ {
						updateType := fmt.Sprintf("updatetype%d", i)
						req := &applet_mgr.CreateStrategy{ProjectName: projectName}
						req.CreateReq.AppletID = appletID
						req.CreateReq.EventType = updateType
						req.CreateReq.Handler = handler

						rsp, _, err := client.CreateStrategy(req)
						NewWithT(t).Expect(err).To(BeNil())
						NewWithT(t).Expect(rsp.EventType).To(Equal(updateType))
					}
				}

				// get list strategies
				{
					req := &applet_mgr.ListStrategy{ProjectName: projectName}
					rsp, _, err := client.ListStrategy(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(num + 1).To(Equal(int(rsp.Total)))
				}

				// remove batch strategies
				{
					req := &applet_mgr.BatchRemoveStrategy{ProjectName: projectName}
					_, err := client.BatchRemoveStrategy(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

			})
		})
	})

	// clear project info
	t.Run("ClearProject", func(t *testing.T) {
		t.Run("#DeleteProject", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// remove project
				{
					req := &applet_mgr.RemoveProject{ProjectName: projectName}
					_, err := client.RemoveProject(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
		})
	})

}
