package mock

import (
	"fmt"
	"net/http"
)

func ExampleResponseWriter() {
	rw := NewMockResponseWriter()
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(`{"status":"ok"}`))

	fmt.Println(string(rw.MustDumpResponse()))
	// Output:
	// HTTP/0.0 200 OK
	// Content-Type: application/json
	//
	// {"status":"ok"}
}
