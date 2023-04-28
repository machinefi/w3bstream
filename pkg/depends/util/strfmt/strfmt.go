package strfmt

//go:generate toolkit gen strfmt -f strfmt.go
const (
	regexpStringProjectName = "^[a-z0-9_]{1,32}$"
)
