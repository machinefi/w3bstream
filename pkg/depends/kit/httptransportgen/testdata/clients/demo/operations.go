// This is a generated source file. DO NOT EDIT
// Source: demo/operations.go

package demo

import (
	"context"
	"mime/multipart"

	"github.com/machinefi/w3bstream/pkg/depends/kit/kit"
	"github.com/machinefi/w3bstream/pkg/depends/kit/metax"
)

type Cookie struct {
	Token string `in:"cookie" name:"token,omitempty"`
}

func (o *Cookie) Path() string {
	return "/demo/cookie"
}

func (o *Cookie) Method() string {
	return "POST"
}

// @StatusErr[ContextCanceled][499000000][ContextCanceled]
// @StatusErr[UnknownError][500000000][UnknownError]

func (o *Cookie) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.Cookie")
	return cli.Do(ctx, o, metas...)
}

func (o *Cookie) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *Cookie) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type Create struct {
	Data Data `in:"body"`
}

func (o *Create) Path() string {
	return "/demo/restful"
}

func (o *Create) Method() string {
	return "POST"
}

func (o *Create) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.Create")
	return cli.Do(ctx, o, metas...)
}

func (o *Create) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*Data, kit.Metadata, error) {
	rsp := new(Data)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *Create) Invoke(cli kit.Client, metas ...kit.Metadata) (*Data, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type DownloadFile struct {
}

func (o *DownloadFile) Path() string {
	return "/demo/binary/files"
}

func (o *DownloadFile) Method() string {
	return "GET"
}

func (o *DownloadFile) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.DownloadFile")
	return cli.Do(ctx, o, metas...)
}

func (o *DownloadFile) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxAttachment, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxAttachment)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *DownloadFile) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxAttachment, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type FormMultipartWithFile struct {
	FormData struct {
		Data  Data                                                                                             `name:"data,omitempty"`
		File  *multipart.FileHeader                                                                            `name:"file"`
		Map   map[GithubComMachinefiW3BstreamPkgDependsKitHttptransportgenTestdataServerPkgTypesProtocol]int32 `name:"map,omitempty"`
		Slice []string                                                                                         `name:"slice,omitempty"`
		// @deprecated
		String string `name:"string,omitempty"`
	} `in:"body" mime:"multipart"`
}

func (o *FormMultipartWithFile) Path() string {
	return "/demo/forms/multipart"
}

func (o *FormMultipartWithFile) Method() string {
	return "POST"
}

func (o *FormMultipartWithFile) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.FormMultipartWithFile")
	return cli.Do(ctx, o, metas...)
}

func (o *FormMultipartWithFile) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *FormMultipartWithFile) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type FormMultipartWithFiles struct {
	FormData struct {
		Files []*multipart.FileHeader `name:"files"`
	} `in:"body" mime:"multipart"`
}

func (o *FormMultipartWithFiles) Path() string {
	return "/demo/forms/multipart-with-files"
}

func (o *FormMultipartWithFiles) Method() string {
	return "POST"
}

func (o *FormMultipartWithFiles) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.FormMultipartWithFiles")
	return cli.Do(ctx, o, metas...)
}

func (o *FormMultipartWithFiles) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *FormMultipartWithFiles) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type FormURLEncoded struct {
	FormData struct {
		Data   Data     `name:"data"`
		Slice  []string `name:"slice"`
		String string   `name:"string"`
	} `in:"body" mime:"urlencoded"`
}

func (o *FormURLEncoded) Path() string {
	return "/demo/forms/urlencoded"
}

func (o *FormURLEncoded) Method() string {
	return "POST"
}

func (o *FormURLEncoded) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.FormURLEncoded")
	return cli.Do(ctx, o, metas...)
}

func (o *FormURLEncoded) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *FormURLEncoded) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type GetByID struct {
	ID       string                                                                                 `in:"path" name:"id" validate:"@string[6,]"`
	Label    []string                                                                               `in:"query" name:"label,omitempty"`
	Name     string                                                                                 `in:"query" name:"name,omitempty"`
	Protocol GithubComMachinefiW3BstreamPkgDependsKitHttptransportgenTestdataServerPkgTypesProtocol `in:"query" name:"protocol,omitempty"`
}

func (o *GetByID) Path() string {
	return "/demo/restful/:id"
}

func (o *GetByID) Method() string {
	return "GET"
}

func (o *GetByID) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.GetByID")
	return cli.Do(ctx, o, metas...)
}

func (o *GetByID) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*Data, kit.Metadata, error) {
	rsp := new(Data)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *GetByID) Invoke(cli kit.Client, metas ...kit.Metadata) (*Data, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type HealthCheck struct {
	PullPolicy GithubComMachinefiW3BstreamPkgDependsKitHttptransportgenTestdataServerPkgTypesPullPolicy `in:"query" name:"pullPolicy,omitempty"`
}

func (o *HealthCheck) Path() string {
	return "/demo/restful"
}

func (o *HealthCheck) Method() string {
	return "HEAD"
}

func (o *HealthCheck) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.HealthCheck")
	return cli.Do(ctx, o, metas...)
}

func (o *HealthCheck) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *HealthCheck) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type Proxy struct {
}

func (o *Proxy) Path() string {
	return "/demo/proxy"
}

func (o *Proxy) Method() string {
	return "GET"
}

// @StatusErr[ClientClosedRequest][499000000][ClientClosedRequest]
// @StatusErr[RequestFailed][500000000][RequestFailed]
// @StatusErr[RequestTransformFailed][400000000][RequestTransformFailed]

func (o *Proxy) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.Proxy")
	return cli.Do(ctx, o, metas...)
}

func (o *Proxy) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*IpInfo, kit.Metadata, error) {
	rsp := new(IpInfo)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *Proxy) Invoke(cli kit.Client, metas ...kit.Metadata) (*IpInfo, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ProxyV2 struct {
}

func (o *ProxyV2) Path() string {
	return "/demo/v2/proxy"
}

func (o *ProxyV2) Method() string {
	return "GET"
}

// @StatusErr[ClientClosedRequest][499000000][ClientClosedRequest]
// @StatusErr[ContextCanceled][499000000][ContextCanceled]
// @StatusErr[RequestFailed][500000000][RequestFailed]
// @StatusErr[RequestTransformFailed][400000000][RequestTransformFailed]
// @StatusErr[UnknownError][500000000][UnknownError]

func (o *ProxyV2) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.ProxyV2")
	return cli.Do(ctx, o, metas...)
}

func (o *ProxyV2) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*IpInfo, kit.Metadata, error) {
	rsp := new(IpInfo)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *ProxyV2) Invoke(cli kit.Client, metas ...kit.Metadata) (*IpInfo, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type Redirect struct {
}

func (o *Redirect) Path() string {
	return "/demo/redirect"
}

func (o *Redirect) Method() string {
	return "GET"
}

func (o *Redirect) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.Redirect")
	return cli.Do(ctx, o, metas...)
}

func (o *Redirect) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *Redirect) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RedirectWhenError struct {
}

func (o *RedirectWhenError) Path() string {
	return "/demo/redirect"
}

func (o *RedirectWhenError) Method() string {
	return "POST"
}

func (o *RedirectWhenError) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.RedirectWhenError")
	return cli.Do(ctx, o, metas...)
}

func (o *RedirectWhenError) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RedirectWhenError) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type RemoveByID struct {
	ID string `in:"path" name:"id" validate:"@string[6,]"`
}

func (o *RemoveByID) Path() string {
	return "/demo/restful/:id"
}

func (o *RemoveByID) Method() string {
	return "DELETE"
}

func (o *RemoveByID) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.RemoveByID")
	return cli.Do(ctx, o, metas...)
}

func (o *RemoveByID) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *RemoveByID) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type ShowImage struct {
}

func (o *ShowImage) Path() string {
	return "/demo/binary/images"
}

func (o *ShowImage) Method() string {
	return "GET"
}

func (o *ShowImage) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.ShowImage")
	return cli.Do(ctx, o, metas...)
}

func (o *ShowImage) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxImagePNG, kit.Metadata, error) {
	rsp := new(GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxImagePNG)
	meta, err := cli.Do(ctx, o, metas...).Into(rsp)
	return rsp, meta, err
}

func (o *ShowImage) Invoke(cli kit.Client, metas ...kit.Metadata) (*GithubComMachinefiW3BstreamPkgDependsKitHttptransportHttpxImagePNG, kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}

type UpdateByID struct {
	ID   string `in:"path" name:"id" validate:"@string[6,]"`
	Data Data   `in:"body"`
}

func (o *UpdateByID) Path() string {
	return "/demo/restful/:id"
}

func (o *UpdateByID) Method() string {
	return "PUT"
}

func (o *UpdateByID) Do(ctx context.Context, cli kit.Client, metas ...kit.Metadata) kit.Result {
	ctx = metax.ContextWith(ctx, "operationID", "demo.UpdateByID")
	return cli.Do(ctx, o, metas...)
}

func (o *UpdateByID) InvokeContext(ctx context.Context, cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	meta, err := cli.Do(ctx, o, metas...).Into(nil)
	return meta, err
}

func (o *UpdateByID) Invoke(cli kit.Client, metas ...kit.Metadata) (kit.Metadata, error) {
	return o.InvokeContext(context.Background(), cli, metas...)
}
