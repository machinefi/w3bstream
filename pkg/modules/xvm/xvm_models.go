package xvm

import "mime/multipart"

type Risc0Info struct {
	ImageId string   `json:"imageId"`
	Elf     string   `json:"elf"`
	Params  []string `json:"params,omitempty"`
}

type CreateRisc0VmReq struct {
	File      *multipart.FileHeader `name:"file"`
	Risc0Info `name:"info"`
}

type CreateRisc0VmRsp struct {
}

type CreateProofReq struct {
	Name         string `json:"name"`
	ImageId      string `json:"imageId"`
	TemplateName string `json:"templateName"`
	InputData    string `json:"inputData"` // json string
}

type CreateProofRsp struct {
	Info string `json:"info"`
}
