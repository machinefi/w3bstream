package risc0vm

type CreateProofReq struct {
	ImageID      string `json:"imageID"`
	PrivateInput string `json:"privateInput"`
	PublicInput  string `json:"publicInput"`
	ReceiptType  string `json:"receiptType"`
}

type CreateProofRsp struct {
	Receipt string `json:"receipt"`
}
