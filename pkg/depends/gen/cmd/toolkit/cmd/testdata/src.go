package testdata

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
	"golang.org/x/mod/sumdb"
)

func F() {
	fmt.Println("demo for code format")
	_, _ = os.Getwd()

	_ = errors.Wrap(errors.New("error string"), "wrap string")

	modfile.ModulePath([]byte(`src`))
	semver.Build("v1.2.3")
	_ = sumdb.Client{}
}
