package strfmt

//go:generate go run ../../kit/validator/strfmt/internal/main/main.go
const (
	regexpStringProjectName = "^[a-z0-9_]{6,32}$"
)
