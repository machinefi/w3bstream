// This is a generated source file. DO NOT EDIT
// Source: testdata/strfmt_generated.go

package testdata

import "github.com/machinefi/w3bstream/pkg/depends/kit/validator"

var AlphaValidator = validator.NewRegexpStrfmtValidator(regexpStringAlpha, "alpha")

func init() {
	validator.DefaultFactory.Register(AlphaValidator)
}
