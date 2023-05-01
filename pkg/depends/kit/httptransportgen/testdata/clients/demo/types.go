// This is a generated source file. DO NOT EDIT
// Source: demo/types.go

package demo

import (
	"bytes"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransportgen/testdata/server/pkg/types"
	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
)

type BytesBuffer = bytes.Buffer

type Data struct {
	ID        string                                                                                 `json:"id"`
	Label     string                                                                                 `json:"label"`
	Protocol  GithubComMachinefiW3BstreamPkgDependsKitHttptransportgenTestdataServerPkgTypesProtocol `json:"protocol,omitempty"`
	PtrString *string                                                                                `json:"ptrString,omitempty"`
	SubData   *SubData                                                                               `json:"subData,omitempty"`
}

type GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxAttachment = httpx.Attachment

type GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxImagePNG = httpx.ImagePNG

type GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxResponse = httpx.Response

type GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxStatusFound struct {
	GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxResponse
}

type GithubComMachinefiW3BstreamPkgDependsKitHttptransportgenTestdataServerPkgTypesProtocol = types.Protocol

type GithubComMachinefiW3BstreamPkgDependsKitHttptransportgenTestdataServerPkgTypesPullPolicy = types.PullPolicy

type GithubComMachinefiW3BstreamPkgDependsKitStatusxErrorField = statusx.ErrorField

type GithubComMachinefiW3BstreamPkgDependsKitStatusxErrorFields = statusx.ErrorFields

type GithubComMachinefiW3BstreamPkgDependsKitStatusxStatusErr = statusx.StatusErr

type IpInfo struct {
	Country     string `json:"country" xml:"country"`
	CountryCode string `json:"countryCode" xml:"countryCode"`
}

type SubData struct {
	Name string `json:"name"`
}
