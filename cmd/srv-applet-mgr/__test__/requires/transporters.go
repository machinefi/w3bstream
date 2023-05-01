package requires

import (
	"net/http"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/global"
	confjwt "github.com/machinefi/w3bstream/pkg/depends/conf/jwt"
	"github.com/machinefi/w3bstream/pkg/depends/kit/sqlx/builder"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/models"
	"github.com/machinefi/w3bstream/pkg/types"
)

func NewAuthPatchRT() func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return &AuthPatchRT{
			tok:  tok,
			next: next,
		}
	}
}

type AuthPatchRT struct {
	tok  string
	next http.RoundTripper
}

func (rt *AuthPatchRT) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+tok)
	return rt.next.RoundTrip(req)
}

var tok string

func init() {
	d := types.MustMgrDBExecutorFromContext(global.Context)
	j := confjwt.MustConfFromContext(global.Context)
	m := &models.Account{}

	err := d.QueryAndScan(
		builder.Select(nil).From(
			d.T(m),
			builder.Where(m.ColRole().Eq(enums.ACCOUNT_ROLE__ADMIN)),
			builder.Limit(1),
		), m,
	)
	if err != nil {
		panic(err)
	}
	tok, err = j.GenerateTokenByPayload(m.AccountID)
	if err != nil {
		panic(err)
	}
}
