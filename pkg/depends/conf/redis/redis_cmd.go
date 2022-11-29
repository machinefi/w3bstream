package redis

func Command(name string, args ...interface{}) *Cmd {
	return &Cmd{
		Name: name,
		Args: args,
	}
}

type Cmd struct {
	Name string
	Args []interface{}
}
