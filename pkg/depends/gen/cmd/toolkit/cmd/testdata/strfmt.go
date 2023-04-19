package testdata

//go:generate toolkit gen strfmt -f strfmt.go
const (
	regexpStringAlpha               = "^[a-zA-Z]+$"
	regexpStringAlphaNumeric        = "^[a-zA-Z0-9]+$"
	regexpStringAlphaUnicode        = "^[\\p{L}]+$"
	regexpStringAlphaUnicodeNumeric = "^[\\p{L}\\p{N}]+$"
	regexpStringNumeric             = "^[-+]?[0-9]+(?:\\.[0-9]+)?$"
	regexpStringNumber              = "^[0-9]+$"
	regexpStringHexadecimal         = "^[0-9a-fA-F]+$"
	regexpStringHexColor            = "^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$"
)
