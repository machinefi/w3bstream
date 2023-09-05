package enums

//go:generate toolkit gen enum TransactionState
type TransactionState uint8

const (
	TRANSACTION_STATE_UNKNOWN TransactionState = iota
	TRANSACTION_STATE__PENDING
	TRANSACTION_STATE__IN_BLOCK
	TRANSACTION_STATE__CONFIRMED
	TRANSACTION_STATE__FAILED
)
