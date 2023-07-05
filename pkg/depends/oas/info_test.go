package oas_test

import (
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func TestInfo(t *testing.T) {
	g := NewCaseGroup("Info")

	g.It("empty", `{"title":"","version":""}`, Info{})

	g.It("with contact", `{"title":"","contact":{"name":"name","url":"url","email":"email"},"version":""}`, Info{
		InfoObject: InfoObject{
			Contact: &Contact{
				ContactObject: ContactObject{
					Name:  "name",
					URL:   "url",
					Email: "email",
				},
			},
		},
	})

	g.It("with licence", `{"title":"","license":{"name":"MIT"},"version":""}`, Info{
		InfoObject: InfoObject{
			License: &License{
				LicenseObject: LicenseObject{
					Name: "MIT",
				},
			},
		},
	})

	g.It("with specification_extensions", `{"title":"","contact":{"x-x":"x"},"license":{"name":"","x-x":"x"},"version":"","x-x":"x"}`, Info{
		InfoObject: InfoObject{
			Contact: &Contact{
				SpecExtensions: SpecExtensions{
					Extensions: map[string]interface{}{
						"x-x": "x",
					},
				},
			},
			License: &License{
				SpecExtensions: SpecExtensions{
					Extensions: map[string]interface{}{
						"x-x": "x",
					},
				},
			},
		},
		SpecExtensions: SpecExtensions{
			Extensions: map[string]interface{}{
				"x-x": "x",
			},
		},
	})

	g.Run(t)
}
