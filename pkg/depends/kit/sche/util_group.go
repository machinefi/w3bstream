package sche

import "sync"

type group struct {
	*sync.WaitGroup
}

func (g *group) Do(f func()) {
	go func() {
		g.Add(1)
		f()
		g.Done()
	}()
}

func (g *group) DoMany(f ...func()) {
	for i := range f {
		g.Do(f[i])
	}
	g.Wait()
}

func Group(wg *sync.WaitGroup) *group {
	return &group{wg}
}

func GroupDo(wg *sync.WaitGroup, f func()) {
	Group(wg).Do(f)
}

func GroupDoMany(f ...func()) {
	(&group{&sync.WaitGroup{}}).DoMany(f...)
}
