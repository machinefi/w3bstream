package roundtrippers

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/propagation"

	"github.com/machinefi/w3bstream/pkg/depends/kit/logr"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
)

type LogRoundTripper struct {
	next http.RoundTripper
}

func NewLogRoundTripper() func(rt http.RoundTripper) http.RoundTripper {
	return func(rt http.RoundTripper) http.RoundTripper {
		return &LogRoundTripper{
			next: rt,
		}
	}
}

func (rt *LogRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)).Inject(ctx, propagation.HeaderCarrier(req.Header))

	ctx, l := logr.Start(ctx, "Request")
	defer l.End()

	cost := timer.Start()
	rsp, err := rt.next.RoundTrip(req.WithContext(ctx))

	defer func() {
		duration := strconv.FormatInt(cost().Microseconds(), 10) + "μs"
		l = l.WithValues(
			"@cst", duration,
			"@mtd", req.Method[0:3],
			"@url", OmitAuthorization(req.URL),
		)

		if err == nil {
			l.Info("success")
		} else {
			l.Warn(errors.Wrap(err, "http request failed"))
		}
	}()

	return rsp, err
}

func OmitAuthorization(u *url.URL) string {
	query := u.Query()
	query.Del("authorization")
	u.RawQuery = query.Encode()
	return u.String()
}
