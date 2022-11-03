package roundtrippers_test

import (
	"net/http"
	"testing"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	. "github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/client/roundtrippers"
)

func TestLogRoundTripper(t *testing.T) {
	mgr := httptransport.NewRequestTsfmFactory(nil, nil)
	mgr.SetDefault()

	req, _ := mgr.NewRequest(http.MethodGet, "https://github.com", nil)
	_, _ = NewLogRoundTripper()(http.DefaultTransport).RoundTrip(req)
}
