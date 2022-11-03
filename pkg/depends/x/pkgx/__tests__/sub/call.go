package sub

func v() int {
	return 2
}

func call() func() int {
	return v
}

type Func func() func() int

func CurryCall() Func {
	if true {
		return func() func() int {
			return func() int {
				return 1
			}
		}
	}
	return func() func() int {
		return call()
	}
}
