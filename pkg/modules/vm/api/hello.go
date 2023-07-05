package api

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("name")
	fmt.Fprintf(w, "hello %v", name)
}
