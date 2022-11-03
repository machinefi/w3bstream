package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/machinefi/w3bstream/pkg/depends/x/misc/timer"
	"github.com/machinefi/w3bstream/pkg/depends/x/pkgx"
)

var Gen = &cobra.Command{
	Use:   "gen",
	Short: "generators",
}

type Generator interface {
	Output(cwd string)
}

func run(cmd string, g func(*pkgx.Pkg) Generator, args ...string) {
	cwd, _ := os.Getwd()
	cost := timer.Start()

	pkg, err := pkgx.LoadFrom(cwd)
	if err != nil {
		panic(err)
	}

	defer func() {
		r := color.New(color.FgHiRed).SprintfFunc()
		g := color.New(color.FgHiGreen).SprintfFunc()
		b := color.New(color.FgHiBlue).SprintfFunc()
		y := color.New(color.FgYellow).SprintfFunc()
		if len(args) == 0 {
			fmt.Printf("%s [%s]\n%s\n",
				r("%-8s", cmd),
				g("%dms", cost().Milliseconds()),
				b("%s", pkg.ID),
			)
		} else {
			fmt.Printf("%s [%s]\n%s %s\n",
				r("%-8s", cmd),
				g("%dms", cost().Milliseconds()),
				b("%s", pkg.ID),
				y("%v", args),
			)
		}
	}()

	g(pkg).Output(cwd)
}
