package oas_test

import (
	"net/url"
	"testing"

	. "github.com/machinefi/w3bstream/pkg/depends/oas"
)

func TestExternalDoc(t *testing.T) {
	g := NewCaseGroup("ExternalDoc")

	g.It("empty", `{}`, ExternalDoc{})

	g.It("with url", `{"url":"https://google.com"}`, ExternalDoc{
		URL: (&url.URL{
			Scheme: "https",
			Host:   "google.com",
		}).String(),
	})

	g.It("with url and description", `{"description":"google","url":"https://google.com"}`, ExternalDoc{
		URL: (&url.URL{
			Scheme: "https",
			Host:   "google.com",
		}).String(),
		Description: "google",
	})

	g.Run(t)
}
