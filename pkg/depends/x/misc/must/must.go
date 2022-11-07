package must

import "log"

// pls add your assert function here, or add type and re-generate

func NoError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func BeTrue(ok bool) {
	if !ok {
		log.Panic("not ok")
	}
}
