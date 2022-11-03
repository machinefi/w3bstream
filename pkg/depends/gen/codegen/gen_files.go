package codegen

import (
	"path/filepath"
	"strings"
)

func IsGoFile(filename string) bool {
	return filepath.Ext(filename) == ".go"
}

func IsGoTestFile(filename string) bool {
	return strings.HasSuffix(filepath.Base(filename), "_test.go")
}

func GenerateFileSuffix(fn string) string {
	dir, base, ext := filepath.Dir(fn), filepath.Base(fn), filepath.Ext(fn)
	if IsGoFile(fn) && IsGoTestFile(fn) {
		base = strings.Replace(base, "_test.go", "__generated_test.go", -1)
	} else {
		base = strings.Replace(base, ext, "__generated"+ext, -1)
	}
	return filepath.Join(dir, base)
}
