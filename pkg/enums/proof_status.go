package enums

//go:generate toolkit gen enum ProofStatus
type ProofStatus uint8

const (
	PROOF_STATUS_UNKNOWN ProofStatus = iota
	PROOF_STATUS__GENERATING
	PROOF_STATUS__SUCCEEDED
	PROOF_STATUS__FAILED
)
