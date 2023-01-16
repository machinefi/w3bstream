package oas_test

import (
	"net/url"
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func TestTag(t *testing.T) {
	g := NewCaseGroup("Tag")

	g.It("empty", `{"name":""}`, Tag{})

	g.It(
		"with external docs",
		`{"name":"tag","externalDocs":{"description":"google","url":"//google.com"}}`,
		func() *Tag {
			t := NewTag("tag")
			t.ExternalDocs = NewExternalDoc((&url.URL{Host: "google.com"}).String(), "google")
			return t
		}(),
	)

	g.Run(t)
}
