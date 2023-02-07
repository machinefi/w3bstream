package routes

import (
	"context"
	"image"
	"image/png"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	testdatapb "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/transformer/testdata"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/testutil/httptransporttestutil/server/pkg/errors"
)

var BinaryRouter = kit.NewRouter(httptransport.Group("/binary"))

func init() {
	RootRouter.Register(BinaryRouter)

	BinaryRouter.Register(kit.NewRouter(DownloadFile{}))
	BinaryRouter.Register(kit.NewRouter(ShowImage{}))
	BinaryRouter.Register(kit.NewRouter(&Protobuf{}))
}

// download file
type DownloadFile struct {
	httpx.MethodGet
}

func (DownloadFile) Path() string {
	return "/files"
}

func (req DownloadFile) Output(ctx context.Context) (interface{}, error) {
	file := httpx.NewAttachment("text.txt", "text/plain")
	file.Write([]byte("123123123"))

	return file, nil
}

// show image
type ShowImage struct {
	httpx.MethodGet
}

func (ShowImage) Path() string {
	return "/images"
}

func (req ShowImage) Output(ctx context.Context) (interface{}, error) {
	i := image.NewAlpha(image.Rectangle{
		Min: image.Pt(0, 0),
		Max: image.Pt(100, 100),
	})

	img := httpx.NewImagePNG()

	if err := png.Encode(img, i); err != nil {
		return nil, err
	}

	return img, nil
}

type Protobuf struct {
	httpx.MethodPost
	testdatapb.Event `in:"body" mime:"protobuf"`
}

func (*Protobuf) Path() string {
	return "/protobuf"
}

func (req *Protobuf) Output(ctx context.Context) (interface{}, error) {
	var (
		evID = req.GetHeader().GetEventId()
		rsp  = &testdatapb.HandleResult{EventId: evID}
	)
	if evID == "" {
		rsp.ErrMsg = "invalid event id"
	} else {
		rsp.Succeeded = true
	}

	pbwr, err := httpx.NewApplicationProtobufWith(rsp)
	if err != nil {
		return nil, errors.InternalServerError
	}
	return pbwr, nil
}
