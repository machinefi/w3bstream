package operators

import (
	"bytes"
	"go/format"
	"text/template"
)

const filterTmplSource = `Filter(func(input interface{}) bool {
	// execute {{ .WasmFuncName }} and cast return result
	return false
})`

type (
	FilterGenerator struct {
		data filterTmplData
	}
	filterTmplData struct {
		WasmFuncName string
	}
)

func NewFilterGenerator(name string) *FilterGenerator {
	return &FilterGenerator{data: filterTmplData{WasmFuncName: name}}
}

func (fg *FilterGenerator) IsValidSuccessor(generator Generator) bool {
	return true
}

func (fg *FilterGenerator) GenCode() (string, error) {
	buffer := new(bytes.Buffer)
	tmpl := template.Must(template.New("").Parse(filterTmplSource))
	if err := tmpl.Execute(buffer, fg.data); err != nil {
		return "", err
	}
	code, err := format.Source(buffer.Bytes())
	if err != nil {
		return "", err
	}
	return string(code), nil
}
