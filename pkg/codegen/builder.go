package codegen

import (
	"errors"

	"github.com/machinefi/w3bstream/pkg/codegen/operators"
)

type ObservableBuilder struct {
	generators []operators.Generator
}

func NewObservableBuilder() *ObservableBuilder {
	return &ObservableBuilder{
		generators: []operators.Generator{&operators.HeadGenerator{}},
	}
}

func (ob *ObservableBuilder) Append(generator operators.Generator) *ObservableBuilder {
	ob.generators = append(ob.generators, generator)

	return ob
}

func (ob *ObservableBuilder) Build() (string, error) {
	retval := ""
	for i := 1; i < len(ob.generators); i++ {
		if !ob.generators[i-1].IsValidSuccessor(ob.generators[i]) {
			return "", errors.New("invalid successor generator")
		}
		code, err := ob.generators[i].GenCode()
		if err != nil {
			return "", err
		}
		retval += "." + code
	}

	return retval, nil
}
