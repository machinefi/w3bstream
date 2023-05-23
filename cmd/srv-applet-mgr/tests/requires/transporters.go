package requires

import (
	"net/http"
	"time"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/pkg/modules/account"
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
	defer Serve()()
	time.Sleep(3 * time.Second)

	_, err := account.CreateAdminIfNotExist(_ctx)
	if err != nil {
		panic(err)
	}
	req := &applet_mgr.LoginByUsername{}
	req.LoginByUsernameReq.Username = "admin"
	req.LoginByUsernameReq.Password = "iotex.W3B.admin"

	rsp, _, err := AuthClient().LoginByUsername(req)
	if err != nil {
		panic(err)
	}
	tok = rsp.Token
}
