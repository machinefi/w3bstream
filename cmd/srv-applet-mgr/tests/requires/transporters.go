package requires

import (
	"net/http"

	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/apis"
	"github.com/machinefi/w3bstream/cmd/srv-applet-mgr/tests/clients/applet_mgr"
	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/retry"
	"github.com/machinefi/w3bstream/pkg/modules/account"
	"github.com/machinefi/w3bstream/pkg/types"
)

func NewAuthPatchRT() func(next http.RoundTripper) http.RoundTripper {
	return NewAuthPatchRTWithToken(tok)
}

func NewAuthPatchRTWithToken(tok string) func(next http.RoundTripper) http.RoundTripper {
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

var (
	tok       string
	AccountID types.SFID
)

func init() {
	go kit.Run(apis.RootMgr, _server.WithContextInjector(_injection))
	go kit.Run(apis.RootEvent, _serverEvent.WithContextInjector(_injection))

	err := retry.Do(retry.Default, func() error {
		if _, _, err := Client().Liveness(); err != nil {
			return err
		}
		if _, _, err := ClientEvent().Liveness(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func init() {
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
	AccountID = rsp.AccountID
}
