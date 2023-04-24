package middleware

import (
	"context"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
	"github.com/machinefi/w3bstream/pkg/enums"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/account_identity"
	"github.com/pkg/errors"
)

func MaybeProjectName(ctx context.Context) (string, bool) {
	if p, ok := ctx.Value("ProjectProvider").(*ProjectProvider); ok && p != nil {
		return p.ProjectName, true
	}
	return "", false
}

func MustProjectName(ctx context.Context) string {
	p, ok := ctx.Value("ProjectProvider").(*ProjectProvider)
	must.BeTrue(ok)
	return p.ProjectName
}

func ProjectNameForDisplay(name string) (string, error) {
	parts := strings.SplitN(name, "_", 3)
	if len(parts) != 3 {
		return "", errors.Errorf("unexpected project name format: %s", name)
	}
	if parts[0] != "aid" && parts[0] != "eth" {
		return "", errors.Errorf("unexpected project name format: %s", name)
	}
	return parts[2], nil
}

func ProjectNameModifier(ctx context.Context) (prefix string, err error) {
	ca, ok := CurrentAccountFromContext(ctx)
	if !ok {
		return "", status.CurrentAccountAbsence
	}

	prefix = "aid_+" + ca.AccountID.String() + "_"
	aci, err := account_identity.GetBySFIDAndType(
		ctx,
		ca.AccountID,
		enums.ACCOUNT_IDENTITY_TYPE__ETHADDRESS,
	)
	if err == nil {
		prefix = "eth_" + aci.IdentityID + "_"
	}
	return
}

// ProjectProvider with account id prefix
type ProjectProvider struct {
	ProjectName string `name:"projectName" in:"path" validate:"@projectName"`
}

func (ProjectProvider) ContextKey() string { return "ProjectProvider" }

func (ProjectProvider) Path() string {
	return "/x/:projectName"
}

func (r *ProjectProvider) Output(ctx context.Context) (interface{}, error) {
	prefix, err := ProjectNameModifier(ctx)
	if err != nil {
		return nil, err
	}
	return &ProjectProvider{
		ProjectName: prefix + r.ProjectName,
	}, nil
}
