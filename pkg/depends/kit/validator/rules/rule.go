package rules

import (
	"bytes"
	"regexp"
)

type Node interface {
	node() // node flag only
	Bytes() []byte
}

type Lit struct {
	Raw []byte
	Node
}

func NewLiteral(raw []byte) *Lit { return &Lit{Raw: raw} }

func (l *Lit) Append(bytes []byte) { l.Raw = append(l.Raw, bytes...) }

func (l *Lit) Bytes() []byte {
	if l == nil {
		return nil
	}
	return l.Raw
}

type Rule struct {
	RAW    []byte
	Name   string
	Params []Node
	Range  []*Lit

	ExclusiveMin bool
	ExclusiveMax bool
	ValueMatrix  [][]*Lit
	Pattern      *regexp.Regexp

	Optional bool
	DftValue []byte

	Node
}

func NewRule(name string) *Rule { return &Rule{Name: name} }

func Parse(rule string) (*Rule, error) { return ParseRaw([]byte(rule)) }

func ParseRaw(b []byte) (*Rule, error) { return NewScanner(b).RootRule() }

func (r *Rule) ComputedValues() []*Lit {
	return ComputedValueMatrix(r.ValueMatrix)
}

func ComputedValueMatrix(matrix [][]*Lit) []*Lit {
	switch len(matrix) {
	case 0:
		return nil
	case 1:
		return matrix[0]
	default:
		ri, li := matrix[0], len(matrix)
		rj, lj := matrix[1], len(matrix)
		values := make([]*Lit, li*lj)
		for i := range ri {
			for j := range rj {
				raw := append(append([]byte{}, ri[i].Bytes()...), rj[j].Bytes()...)
				values[i*lj+j] = NewLiteral(raw)
			}
		}
		return ComputedValueMatrix(append([][]*Lit{values}, matrix[2:]...))
	}
}

func (r *Rule) Bytes() []byte {
	if r == nil {
		return nil
	}

	buf := &bytes.Buffer{}
	buf.WriteByte('@')
	buf.WriteString(r.Name)

	if len(r.Params) > 0 {
		buf.WriteByte('<')
		for i, p := range r.Params {
			if i > 0 {
				buf.WriteByte(',')
			}
			if p != nil {
				buf.Write(p.Bytes())
			}
		}
		buf.WriteByte('>')
	}

	if len(r.Range) > 0 {
		if r.ExclusiveMin {
			buf.WriteRune('(')
		} else {
			buf.WriteRune('[')
		}
		for i, p := range r.Range {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.Write(p.Bytes())
		}
		if r.ExclusiveMax {
			buf.WriteRune(')')
		} else {
			buf.WriteRune(']')
		}
	}

	for i := range r.ValueMatrix {
		values := r.ValueMatrix[i]

		buf.WriteByte('{')

		for i, p := range values {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.Write(p.Bytes())
		}

		buf.WriteByte('}')
	}

	if r.Pattern != nil {
		buf.Write(Slash([]byte(r.Pattern.String())))
	}

	if r.Optional {
		if r.DftValue != nil {
			buf.WriteByte(' ')
			buf.WriteByte('=')
			buf.WriteByte(' ')

			buf.Write(SingleQuote(r.DftValue))
		} else {
			buf.WriteByte('?')
		}
	}

	return buf.Bytes()
}
