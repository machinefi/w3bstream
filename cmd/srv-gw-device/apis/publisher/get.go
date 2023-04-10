package publisher

import "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"

// TODO migration srv-applet-mgr/v0/publisher

type GetPublisherByID struct {
	httpx.MethodGet
}

type ListPublisher struct {
}

type CreatePublisher struct{}

type CreatePublisherByDID struct{}
