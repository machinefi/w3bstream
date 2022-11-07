package main

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

var patcher = []byte(`package runtime

func GoID() int64 {
    return getg().goid
}
`)

func main() {
	pkg, err := build.Default.Import("runtime", "", build.FindOnly)
	if err != nil {
		fmt.Println("err:", err)
	}
	err = ioutil.WriteFile(path.Join(pkg.Dir, "proc_id.go"), patcher, os.ModePerm)
	if err != nil {
		fmt.Println("err:", err)
	}
	_, err = exec.Command("go", "install", "runtime").CombinedOutput()
	if err != nil {
		fmt.Println("err:", err)
	}
}
