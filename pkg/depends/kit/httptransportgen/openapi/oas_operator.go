package openapi

import (
	"go/types"
	"net/http"
	"sort"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
	"github.com/machinefi/w3bstream/pkg/depends/oas"
)

type Operator struct {
	httptransport.RouteMeta

	Tag         string
	Description string

	NonBodyParams map[string]*oas.Parameter
	RequestBody   *oas.RequestBody

	StatusErrors      []*statusx.StatusErr
	StatusErrorSchema *oas.Schema

	SuccessStatus   int
	SuccessType     types.Type
	SuccessResponse *oas.Response
}

func (o *Operator) AddNonBodyParameter(param *oas.Parameter) {
	if o.NonBodyParams == nil {
		o.NonBodyParams = map[string]*oas.Parameter{}
	}
	o.NonBodyParams[param.Name] = param
}

func (o *Operator) SetRequestBody(body *oas.RequestBody) {
	o.RequestBody = body
}

func (o *Operator) BindOperation(mtd string, opt *oas.Operation, last bool) {
	// parameters
	params := map[string]bool{}
	for _, param := range opt.Parameters {
		params[param.Name] = true
	}

	for _, parameter := range o.NonBodyParams {
		if !params[parameter.Name] {
			opt.Parameters = append(opt.Parameters, parameter)
		}
	}

	// request body
	if o.RequestBody != nil {
		opt.SetRequestBody(o.RequestBody)
	}

	// status errors
	for _, se := range o.StatusErrors {
		statuserrs := make([]string, 0)

		if opt.Responses.Responses != nil {
			if rsp, ok := opt.Responses.Responses[se.StatusCode()]; ok {
				if rsp.Extensions != nil {
					if v, ok := rsp.Extensions[XStatusErrs]; ok {
						if list, ok := v.([]string); ok {
							statuserrs = append(statuserrs, list...)
						}
					}
				}
			}
		}
		statuserrs = append(statuserrs, se.Summary())

		sort.Strings(statuserrs)

		rsp := oas.NewResponse("")
		rsp.AddExtension(XStatusErrs, statuserrs)
		rsp.AddContent(httpx.MIME_JSON, oas.NewMediaTypeWithSchema(o.StatusErrorSchema))
		opt.AddResponse(se.StatusCode(), rsp)
	}

	// last bind response
	if last {
		opt.OperationId = o.ID
		opt.Deprecated = o.Deprecated
		opt.Summary = o.Summary
		opt.Description = o.Description

		if o.Tag != "" {
			opt.Tags = []string{o.Tag}
		}

		if o.SuccessType == nil {
			opt.Responses.AddResponse(http.StatusNoContent, &oas.Response{})
		} else {
			status := o.SuccessStatus
			if status == 0 {
				status = http.StatusOK
				if mtd == http.MethodPost {
					status = http.StatusCreated
				}
			}
			if status >= http.StatusMultipleChoices && status < http.StatusBadRequest {
				o.SuccessResponse = oas.NewResponse(o.SuccessResponse.Description)
			}
			opt.Responses.AddResponse(status, o.SuccessResponse)
		}
	}

	// sort all parameters by position and name
	if len(opt.Parameters) > 0 {
		sort.Slice(opt.Parameters, func(i, j int) bool {
			return ParamPosOrder[opt.Parameters[i].In]+opt.Parameters[i].Name <
				ParamPosOrder[opt.Parameters[j].In]+opt.Parameters[j].Name
		})
	}
}

type OperatorWithTypeName struct {
	*Operator
	TypeName *types.TypeName
}

func (op *OperatorWithTypeName) String() string {
	return op.TypeName.Pkg().Name() + "." + op.TypeName.Name()
}
