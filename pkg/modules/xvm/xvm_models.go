package xvm

import "mime/multipart"

type Risc0Info struct {
	ImageId string   `json:"imageId"`
	Elf     string   `json:"elf"`
	Params  []string `json:"params"`
}

type CreateRisc0VmReq struct {
	File      *multipart.FileHeader `name:"file"`
	Risc0Info `name:"info"`
}

type CreateRisc0VmRsp struct {
}
