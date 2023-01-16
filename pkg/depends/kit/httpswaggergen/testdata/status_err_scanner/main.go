package main

import (
	"fmt"
	"net/http"

	"github.com/machinefi/w3bstream/pkg/depends/kit/statusx"
)

// @StatusErr[InternalServerError][500100001][InternalServerError]
func call() {
	fn()
}

func main() {
	call()
	fmt.Println(Unauthorized)
}

func fn() error {
	return statusx.Wrap(fmt.Errorf("test"), http.StatusInternalServerError, "Test")
}
