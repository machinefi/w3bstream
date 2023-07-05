package oas_test

import (
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func TestSpecExtensions(t *testing.T) {
	g := NewCaseGroup("SpecExtensions")

	g.It("empty", `{}`, SpecExtensions{})

	g.It("with extensions", `{"x-a":"xxx"}`, func() *SpecExtensions {
		e := &SpecExtensions{}
		e.AddExtension("x-b", nil)
		e.AddExtension("x-a", "xxx")
		return e
	}())

	g.Run(t)
}
