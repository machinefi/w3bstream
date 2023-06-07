package integrations

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/requires"
	confid "github.com/machinefi/w3bstream/pkg/depends/conf/id"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
	"github.com/machinefi/w3bstream/pkg/types"
)

func TestPublisherAPIs(t *testing.T) {
	var (
		ctx           = requires.Context()
		client        = requires.AuthClient()
		projectName   = "test_publisher_project"
		publisherName = "testpublisher"
		publisherKey  = confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID().String()

		publisherID types.SFID
	)

	t.Logf("random a project name: %s, use this name create a project.", projectName)

	{
		req := &applet_mgr.CreateProject{}
		req.CreateReq.Name = projectName

		_, _, err := client.CreateProject(req)
		if err != nil {
			panic(err)
		}
	}

	defer func() {
		req := &applet_mgr.RemoveProject{ProjectName: projectName}
		_, err := client.RemoveProject(req)
		if err != nil {
			panic(err)
		}
	}()

	t.Logf("random a publisher name and publisehr key: %s - %s, then create a pulbisher .",
		publisherName, publisherKey)

	t.Run("Publisher", func(t *testing.T) {
		t.Run("#CreatePublisher", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// create publisher
				{
					req := &applet_mgr.CreatePublisher{
						ProjectName: projectName,
						CreateReq: publisher.CreateReq{
							Name: publisherName,
							Key:  publisherKey,
						},
					}

					rsp, _, err := client.CreatePublisher(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.Name).To(Equal(publisherName))
					publisherID = rsp.PublisherID
				}

				// get publisher
				{
					req := &applet_mgr.GetPublisher{PublisherID: publisherID}
					rsp, _, err := client.GetPublisher(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(rsp.Name).To(Equal(publisherName))
				}

				// update publisher
				{
					updateName := "updatepublisher"
					req := &applet_mgr.UpdatePublisher{
						ProjectName: projectName,
						PublisherID: publisherID,
						UpdateReq: publisher.UpdateReq{
							Name: updateName,
							Key:  publisherKey,
						},
					}
					_, err := client.UpdatePublisher(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

				// remove publisher
				{
					req := &applet_mgr.RemovePublisher{PublisherID: publisherID}
					_, err := client.RemovePublisher(req)
					NewWithT(t).Expect(err).To(BeNil())
				}
			})
		})
	})

	t.Run("BatchPublisher", func(t *testing.T) {
		t.Run("#CreatePublishers", func(t *testing.T) {
			t.Run("#Success", func(t *testing.T) {

				// prepare data
				num := 5
				{
					for i := 0; i < num; i++ {
						pubName := fmt.Sprintf("testpublisher%d", i)
						req := &applet_mgr.CreatePublisher{
							ProjectName: projectName,
							CreateReq: publisher.CreateReq{
								Name: pubName,
								Key:  confid.MustSFIDGeneratorFromContext(ctx).MustGenSFID().String(),
							},
						}
						rsp, _, err := client.CreatePublisher(req)
						NewWithT(t).Expect(err).To(BeNil())
						NewWithT(t).Expect(rsp.Name).To(Equal(pubName))
					}
				}

				// get list publisher
				{
					req := &applet_mgr.ListPublisher{ProjectName: projectName}
					rsp, _, err := client.ListPublisher(req)
					NewWithT(t).Expect(err).To(BeNil())
					NewWithT(t).Expect(num).To(Equal(int(rsp.Total)))
				}

				// remove batch publisher
				{
					req := &applet_mgr.BatchRemovePublisher{ProjectName: projectName}
					_, err := client.BatchRemovePublisher(req)
					NewWithT(t).Expect(err).To(BeNil())
				}

			})
		})
	})
}
