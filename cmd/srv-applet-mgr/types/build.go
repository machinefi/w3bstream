package types

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/machinefi/w3bstream/pkg/depends/base/consts"
)

var (
	Name      string
	Feature   string
	Version   string
	Timestamp string
	Group     string

	BuildVersion string
)

func init() {
	if Name == "" {
		Name = "srv-applet-mgr"
	}
	if Feature == "" {
		Feature = "unknown"
	}
	if Version == "" {
		Version = "unknown"
	}
	if Timestamp == "" {
		Timestamp = "unknown"
	}
	if Group == "" {
		Group = "srv-applet-mgr"
	}
	_ = os.Setenv(consts.EnvProjectName, Name)
	_ = os.Setenv(consts.EnvProjectFeat, Feature)
	_ = os.Setenv(consts.EnvProjectVersion, Version)
	_ = os.Setenv(consts.EnvResourceGroup, Group)

	BuildVersion = fmt.Sprintf("%s@%s_%s", Feature, Version, Timestamp)

	fmt.Printf(color.CyanString("%s: %s\n\n", Name, BuildVersion))
}
