package tasks

import "github.com/iotexproject/Bumblebee/kit/kit"

var Root = kit.NewRouter()

func init() {
	Root.Register(kit.NewRouter(&HandleEvent{}))
}
