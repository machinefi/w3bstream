package did

import (
	"bytes"
	"context"
	"net/http"

	"github.com/machinefi/w3bstream/pkg/depends/x/contextx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

type DID struct {
	VCVerificationService string `env:""`
}

func (d *DID) SetDefault() {}

func (d *DID) Init() {}

func (d *DID) CheckVC(vc []byte) (bool, error) {
	req, err := http.NewRequest("GET", d.VCVerificationService, bytes.NewBuffer(vc))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/did+ld+json")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, nil
}

type did struct{}

func WithDIDContext(d *DID) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, did{}, d)
	}
}

func DIDFromContext(ctx context.Context) (*DID, bool) {
	d, ok := ctx.Value(did{}).(*DID)
	return d, ok
}

func MustDIDFromContext(ctx context.Context) *DID {
	d, ok := ctx.Value(did{}).(*DID)
	must.BeTrue(ok)
	return d
}
