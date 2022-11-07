// This is a generated source file. DO NOT EDIT
// Source: httpx_test/httpx_redirect__generated_test.go

package httpx_test

import (
	"fmt"
	"net/url"

	"github.com/machinefi/w3bstream/pkg/depends/kit/httptransport/httpx"
)

func ExampleStatusMultipleChoices() {
	m := httpx.RedirectWithStatusMultipleChoices(&(url.URL{
		Path: "/test",
	}))
	fmt.Println(m.StatusCode())
	fmt.Println(m.Location())
	// Output:
	// 300
	// /test
}

func ExampleStatusMovedPermanently() {
	m := httpx.RedirectWithStatusMovedPermanently(&(url.URL{
		Path: "/test",
	}))
	fmt.Println(m.StatusCode())
	fmt.Println(m.Location())
	// Output:
	// 301
	// /test
}

func ExampleStatusFound() {
	m := httpx.RedirectWithStatusFound(&(url.URL{
		Path: "/test",
	}))
	fmt.Println(m.StatusCode())
	fmt.Println(m.Location())
	// Output:
	// 302
	// /test
}

func ExampleStatusSeeOther() {
	m := httpx.RedirectWithStatusSeeOther(&(url.URL{
		Path: "/test",
	}))
	fmt.Println(m.StatusCode())
	fmt.Println(m.Location())
	// Output:
	// 303
	// /test
}

func ExampleStatusNotModified() {
	m := httpx.RedirectWithStatusNotModified(&(url.URL{
		Path: "/test",
	}))
	fmt.Println(m.StatusCode())
	fmt.Println(m.Location())
	// Output:
	// 304
	// /test
}

func ExampleStatusUseProxy() {
	m := httpx.RedirectWithStatusUseProxy(&(url.URL{
		Path: "/test",
	}))
	fmt.Println(m.StatusCode())
	fmt.Println(m.Location())
	// Output:
	// 305
	// /test
}

func ExampleStatusTemporaryRedirect() {
	m := httpx.RedirectWithStatusTemporaryRedirect(&(url.URL{
		Path: "/test",
	}))
	fmt.Println(m.StatusCode())
	fmt.Println(m.Location())
	// Output:
	// 307
	// /test
}

func ExampleStatusPermanentRedirect() {
	m := httpx.RedirectWithStatusPermanentRedirect(&(url.URL{
		Path: "/test",
	}))
	fmt.Println(m.StatusCode())
	fmt.Println(m.Location())
	// Output:
	// 308
	// /test
}
