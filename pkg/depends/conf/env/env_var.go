package env

type Var struct {
	Name  string
	Value string
	Mask  string

	Optional      bool
	IsUpstream    bool
	IsCopy        bool
	IsExpose      bool
	IsHealthCheck bool
}

func (v *Var) Key(prefix string) string {
	if prefix != "" {
		return prefix + "__" + v.Name
	}
	return v.Name
}

func (v *Var) SetMeta(flags map[string]bool) {
	// TODO config meta
}
