// This is a generated source file. DO NOT EDIT
// Source: strfmt/strfmt_generated.go

package strfmt

import "github.com/machinefi/w3bstream/pkg/depends/kit/validator"

func init() {
	validator.DefaultFactory.Register(ProjectNameValidator)
}

var ProjectNameValidator = validator.NewRegexpStrfmtValidator(regexpStringProjectName, "project-name", "projectName")
