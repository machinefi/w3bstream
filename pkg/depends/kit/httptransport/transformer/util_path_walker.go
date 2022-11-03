package transformer

import (
	"fmt"
	"strings"
)

type PathWalker struct {
	path []interface{}
}

func (p *PathWalker) Enter(i interface{}) { p.path = append(p.path, i) }

func (p *PathWalker) Exit() { p.path = p.path[:len(p.path)-1] }

func (p *PathWalker) Paths() []interface{} { return p.path }

func (p *PathWalker) String() string {
	b := &strings.Builder{}
	for i := 0; i < len(p.path); i++ {
		switch x := p.path[i].(type) {
		case string:
			if b.Len() != 0 {
				b.WriteByte('.')
			}
			b.WriteString(x)
		case int:
			b.WriteString(fmt.Sprintf("[%d]", x))
		}
	}
	return b.String()
}
