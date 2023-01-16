// This is a generated source file. DO NOT EDIT
// Source: testdata/strfmt_generated.go

package testdata

import "github.com/machinefi/w3bstream/pkg/depends/kit/validator"

var (
	AlphaValidator               = validator.NewRegexpStrfmtValidator(regexpStringAlpha, "alpha")
	AlphaNumericValidator        = validator.NewRegexpStrfmtValidator(regexpStringAlphaNumeric, "alpha-numeric", "alphaNumeric")
	AlphaUnicodeValidator        = validator.NewRegexpStrfmtValidator(regexpStringAlphaUnicode, "alpha-unicode", "alphaUnicode")
	AlphaUnicodeNumericValidator = validator.NewRegexpStrfmtValidator(regexpStringAlphaUnicodeNumeric, "alpha-unicode-numeric", "alphaUnicodeNumeric")
	HexColorValidator            = validator.NewRegexpStrfmtValidator(regexpStringHexColor, "hex-color", "hexColor")
	HexadecimalValidator         = validator.NewRegexpStrfmtValidator(regexpStringHexadecimal, "hexadecimal")
	NumberValidator              = validator.NewRegexpStrfmtValidator(regexpStringNumber, "number")
	NumericValidator             = validator.NewRegexpStrfmtValidator(regexpStringNumeric, "numeric")
)

func init() {
	validator.DefaultFactory.Register(AlphaValidator)
	validator.DefaultFactory.Register(AlphaNumericValidator)
	validator.DefaultFactory.Register(AlphaUnicodeValidator)
	validator.DefaultFactory.Register(AlphaUnicodeNumericValidator)
	validator.DefaultFactory.Register(HexColorValidator)
	validator.DefaultFactory.Register(HexadecimalValidator)
	validator.DefaultFactory.Register(NumberValidator)
	validator.DefaultFactory.Register(NumericValidator)
}
