package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"

	"github.com/machinefi/w3bstream/pkg/depends/gen/codegen/formatx"
	"github.com/machinefi/w3bstream/pkg/depends/x/misc/must"
)

var (
	Format *cobra.Command
	grp    string
	dir    string
)

func init() {
	cmd := &cobra.Command{
		Use:   "fmt",
		Short: "format go code",
		Run: func(cmd *cobra.Command, args []string) {
			PrepareArgs()
			MustFormatRoot(grp, dir)
		},
	}

	cmd.PersistentFlags().StringVarP(&grp, "group", "g", "", "grouped imports, default find in go.mod in current dir")
	cmd.PersistentFlags().StringVarP(&dir, "dir", "d", ".", "format directory")

	Format = cmd
}

func ModulePath(path string) (string, error) {
	dat, err := os.ReadFile(path)
	if err == nil {
		return modfile.ModulePath(dat), nil
	}

	return "", err
}

func PrepareArgs() {
	cwd, _ := os.Getwd()

	if grp == "" {
		grp, _ = ModulePath(filepath.Join(cwd, "go.mod"))
	}
	if dir == "." {
		dir = cwd
	} else {
		dir = filepath.Join(cwd, dir)
	}

	fmt.Printf("module name: %s\n", grp)
	fmt.Printf("source root: %s\n", dir)
}

func MustFormatRoot(group, root string) { must.NoError(FormatRoot(group, root, false)) }

func FormatRoot(group, root string, verbose bool) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		code, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		formatted, err := formatx.Format(path, group, code, formatx.SortImports)
		if err != nil {
			return err
		}
		if bytes.Equal(code, formatted) {
			return nil
		}
		fmt.Println(path)
		if verbose {
		}
		return os.WriteFile(path, formatted, stat.Mode())
	})
}
