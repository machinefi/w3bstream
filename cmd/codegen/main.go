package main

import (
	"github.com/machinefi/w3bstream/pkg/codegen"
	"github.com/machinefi/w3bstream/pkg/depends/conf/log"
)

func main() {
	if err := codegen.NewCodeGenCmd().Execute(); err != nil {
		log.Std().Fatal(err)
	}
}
