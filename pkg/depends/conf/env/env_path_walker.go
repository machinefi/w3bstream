package env

import (
	"bytes"
	"strconv"
)

type PathWalker struct{ path []interface{} }

func NewPathWalker() *PathWalker { return &PathWalker{} }

func (p *PathWalker) Enter(i interface{}) { p.path = append(p.path, i) }

func (p *PathWalker) Exit() { p.path = p.path[:len(p.path)-1] }

func (p *PathWalker) Paths() []interface{} { return p.path }

func (p *PathWalker) String() string { return StringifyPath(p.path...) }

func StringifyPath(paths ...interface{}) string {
	buf := bytes.NewBuffer(nil)
	for i, key := range paths {
		if i > 0 {
			buf.WriteRune('_')
		}
		switch v := key.(type) {
		case string:
			buf.WriteString(v)
		case int:
			buf.WriteString(strconv.Itoa(v))
		}
	}
	return buf.String()
}
