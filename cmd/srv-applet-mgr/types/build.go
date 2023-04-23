package types

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
)

var (
	Name      string
	Remote    string
	Branch    string
	Commit    string
	Version   string
	Timestamp string
)

func init() {
	_ = os.Setenv(consts.EnvProjectName, Name)
	_ = os.Setenv(consts.EnvProjectFeat, Branch+"@"+Commit)
	_ = os.Setenv(consts.EnvProjectVersion, Version)

	fmt.Printf(color.CyanString(
		"\n%s:%s was built at %s on %s(%s)\n\n",
		Remote, Name, Timestamp, Branch, Commit,
	))
}
