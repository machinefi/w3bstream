package risc0vm

type CreateProofReq struct {
	ImageID string   `json:"imageID"`
	Params  []string `json:"params"`
}

type CreateProofRsp struct {
	Receipt string `json:"receipt"`
}
