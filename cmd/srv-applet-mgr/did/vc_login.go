package did

import (
	"context"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/conf/did"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/errors/status"
	"github.com/machinefi/w3bstream/pkg/modules/project"
	"github.com/machinefi/w3bstream/pkg/modules/publisher"
)

// TODO complete definition
type VC struct {
	ID string `json:"id"`
}

type LoginRsp struct {
	Token string `json:"token"`
}

type VCLogin struct {
	httpx.MethodPost
	ProjectName string `in:"path" name:"projectName"`
}

func (r *VCLogin) Path() string {
	return "/:projectName"
}

func (r *VCLogin) Output(ctx context.Context) (interface{}, error) {
	prj, err := project.GetProjectByProjectName(ctx, r.ProjectName)
	if err != nil {
		return nil, err
	}
	vcByte := []byte{}
	vc := VC{} // TODO read vc
	did := did.MustDIDFromContext(ctx)
	ok, err := did.CheckVC(vcByte)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.BadRequest.StatusErr().WithDesc("verifiable credentials illegal")
	}

	req := publisher.CreatePublisherReq{
		Key: vc.ID,
	}
	p, err := publisher.CreatePublisher(ctx, prj.ProjectID, &req)
	if err != nil {
		if !strings.Contains(err.Error(), "duplicate key") {
			return nil, err
		}
		p, err = publisher.GetPublisherByPubKeyAndProjectName(ctx, vc.ID, r.ProjectName)
		if err != nil {
			return nil, err
		}
	}
	return &LoginRsp{Token: p.Token}, nil
}
