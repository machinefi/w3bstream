package ver

type Versions []Version

func (c Versions) Len() int {
	return len(c)
}

func (c Versions) Less(i, j int) bool {
	return c[i].LessThan(&c[j])
}

func (c Versions) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
