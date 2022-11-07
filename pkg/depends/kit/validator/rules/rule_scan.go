package rules

import (
	"bytes"
	"regexp"
	textscanner "text/scanner"

	"github.com/machinefi/w3bstream/pkg/depends/kit/validator/errors"
)

type scanner struct {
	raw []byte
	*textscanner.Scanner
}

func NewScanner(b []byte) *scanner {
	s := &textscanner.Scanner{}
	s.Init(bytes.NewReader(b))
	return &scanner{b, s}
}

func (s *scanner) RootRule() (*Rule, error) {
	rule, err := s.rule()
	if err != nil {
		return nil, err
	}
	if tok := s.Scan(); tok != EOF {
		return nil, errors.NewSyntaxError(
			"%s | rule should be end but got `%s`",
			s.raw[0:s.Pos().Offset], string(tok))
	}
	return rule, nil
}

func (s *scanner) rule() (*Rule, error) {
	// simple          @name
	// with parameters @name<param> @name<param1,param2...>
	// with ranges     @name[from,to), @name[length]
	// with values     @name{value1,value2}
	// with regexp     @name\/d+/
	// optional        @name?
	// default value   @name=value @name='xxx'
	// compose         @map<@string[1,10],@string{A,B,C}>[0,10]
	if first := s.Next(); first != '@' {
		return nil, errors.NewSyntaxError(
			"%s | rule should start with `@` but got `%s`",
			s.raw[0:s.Pos().Offset], string(first),
		)
	}
	start := s.Pos().Offset - 1
	name, err := s.lit()
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errors.NewSyntaxError(
			"%s | rule missing name",
			s.raw[0:s.Pos().Offset],
		)
	}
	r := NewRule(name)
LOOP:
	for tok := s.Peek(); ; tok = s.Peek() {
		switch tok {
		default:
			break LOOP
		case ' ':
			s.Next()
		case '?', '=':
			optional, dftv, err := s.inherent()
			if err != nil {
				return nil, err
			}
			r.Optional, r.DftValue = optional, dftv
		case '<':
			params, err := s.params()
			if err != nil {
				return nil, err
			}
			r.Params = params
		case '[', '(':
			ranges, end, err := s.ranges()
			if err != nil {
				return nil, err
			}
			r.Range = ranges
			r.ExclusiveMin = tok == '('
			r.ExclusiveMax = end == ')'
		case '{':
			values, err := s.values()
			if err != nil {
				return nil, err
			}
			r.ValueMatrix = append(r.ValueMatrix, values)
		case '/':
			pattern, err := s.pattern()
			if err != nil {
				return nil, err
			}
			r.Pattern = pattern
		}
	}

	end := s.Pos().Offset
	r.RAW = s.raw[start:end]
	return r, nil
}

func (s *scanner) lit() (string, error) {
	tok := s.Scan()
	if keychars[tok] {
		return "", errors.NewSyntaxError(
			"%s | invalid literal token `%s`",
			s.raw[0:s.Pos().Offset], string(tok),
		)
	}
	return s.TokenText(), nil
}

// inherent optional or default value
func (s *scanner) inherent() (bool, []byte, error) {
	first := s.Next()
	if !(first == '=' || first == '?') {
		return false, nil, errors.NewSyntaxError(
			"%s | optional or default value of rule should start with `?` or `=`",
			s.raw[0:s.Pos().Offset],
		)
	}

	b := &bytes.Buffer{}

	tok := s.Peek()
	for tok == ' ' {
		tok = s.Next()
	}

	if tok == '\'' {
		for tok = s.Peek(); tok != '\''; tok = s.Peek() {
			if tok == EOF {
				return true, nil, errors.NewSyntaxError(
					"%s | default value of of rule should end with `'`",
					s.raw[0:s.Pos().Offset],
				)
			}
			if tok == '\\' {
				tok = s.Next()
				next := s.Next()
				// \' -> '
				if next != '\'' {
					b.WriteRune(tok)
				}
				b.WriteRune(next)
				continue
			}
			b.WriteRune(tok)
			s.Next()
		}
		s.Next()
	} else if tok != EOF && tok != '>' && tok != ',' {
		// end or in stmt
		b.WriteRune(tok)
		lit, err := s.lit()
		if err != nil {
			return false, nil, err
		}
		b.WriteString(lit)
	}

	dftv := b.Bytes()

	if first == '=' && dftv == nil {
		return true, []byte{}, nil
	}

	return true, dftv, nil
}

func (s *scanner) params() ([]Node, error) {
	if first := s.Next(); first != '<' {
		return nil, errors.NewSyntaxError(
			"%s | parameters of rule should start with `<` but got `%s`",
			s.raw[0:s.Pos().Offset], string(first),
		)
	}

	params := map[int]Node{}
	paramc := 1

	for tok := s.Peek(); tok != '>'; tok = s.Peek() {
		if tok == EOF {
			return nil, errors.NewSyntaxError(
				"%s | parameters of rule should end with `>` but got `%s`",
				s.raw[0:s.Pos().Offset], string(tok),
			)
		}
		switch tok {
		case ' ':
			s.Next()
		case ',':
			s.Next()
			paramc++
		case '@':
			rule, err := s.rule()
			if err != nil {
				return nil, err
			}
			params[paramc] = rule
		default:
			raw, err := s.lit()
			if err != nil {
				return nil, err
			}
			if node, ok := params[paramc]; !ok {
				params[paramc] = NewLiteral([]byte(raw))
			} else if lit, ok := node.(*Lit); ok {
				lit.Append([]byte(raw))
			} else {
				return nil, errors.NewSyntaxError(
					"%s | rule should be end but got `%s`",
					s.raw[0:s.Pos().Offset], string(tok),
				)
			}
		}
	}

	lst := make([]Node, paramc)
	for i := range lst {
		if p, ok := params[i+1]; ok {
			lst[i] = p
		} else {
			lst[i] = NewLiteral([]byte(""))
		}
	}

	s.Next()
	return lst, nil
}

func (s *scanner) ranges() ([]*Lit, rune, error) {
	if first := s.Next(); !(first == '[' || first == '(') {
		return nil, first, errors.NewSyntaxError(
			"%s range of rule should start with `[` or `(` but got `%s`",
			s.raw[0:s.Pos().Offset], string(first),
		)
	}

	lits := map[int]*Lit{}
	litc := 1

	for tok := s.Peek(); !(tok == ']' || tok == ')'); tok = s.Peek() {
		if tok == EOF {
			return nil, tok, errors.NewSyntaxError(
				"%s range of rule should end with `]` `)` but got `%s`",
				s.raw[0:s.Pos().Offset], string(tok),
			)
		}
		switch tok {
		case ' ':
			s.Next()
		case ',':
			s.Next()
			litc++
		default:
			raw, err := s.lit()
			if err != nil {
				return nil, tok, err
			}
			if lit, ok := lits[litc]; !ok {
				lits[litc] = NewLiteral([]byte(raw))
			} else {
				lit.Append([]byte(raw))
			}
		}
	}

	lst := make([]*Lit, litc)

	for i := range lst {
		if p, ok := lits[i+1]; ok {
			lst[i] = p
		} else {
			lst[i] = NewLiteral([]byte(""))
		}
	}

	return lst, s.Next(), nil
}

func (s *scanner) values() ([]*Lit, error) {
	if first := s.Next(); first != '{' {
		return nil, errors.NewSyntaxError(
			"%s | vals of rule should start with `{` but got `%s`",
			s.raw[0:s.Pos().Offset], string(first))
	}

	vals := map[int]*Lit{}
	valc := 1

	for tok := s.Peek(); tok != '}'; tok = s.Peek() {
		if tok == EOF {
			return nil, errors.NewSyntaxError(
				"%s vals of rule should end with `}`",
				s.raw[0:s.Pos().Offset],
			)
		}
		switch tok {
		case ' ':
			s.Next()
		case ',':
			s.Next()
			valc++
		default:
			raw, err := s.lit()
			if err != nil {
				return nil, err
			}
			if literal, ok := vals[valc]; !ok {
				vals[valc] = NewLiteral([]byte(raw))
			} else {
				literal.Append([]byte(raw))
			}
		}
	}
	s.Next()

	lst := make([]*Lit, valc)
	for i := range lst {
		if p, ok := vals[i+1]; ok {
			lst[i] = p
		} else {
			lst[i] = NewLiteral([]byte(""))
		}
	}
	return lst, nil
}

func (s *scanner) pattern() (*regexp.Regexp, error) {
	if first := s.Next(); first != '/' {
		return nil, errors.NewSyntaxError(
			"%s | pattern of rule should start with `/`",
			s.raw[0:s.Pos().Offset],
		)
	}

	b := &bytes.Buffer{}

	for tok := s.Peek(); tok != '/'; tok = s.Peek() {
		if tok == EOF {
			return nil, errors.NewSyntaxError(
				"%s | pattern of rule should end with `/`",
				s.raw[0:s.Pos().Offset],
			)
		}
		if tok == '\\' {
			tok = s.Next()
			next := s.Next()
			// \/ -> /
			if next != '/' {
				b.WriteRune(tok)
			}
			b.WriteRune(next)
			continue
		}
		b.WriteRune(tok)
		s.Next()
	}
	s.Next()

	return regexp.Compile(b.String())
}

var keychars = map[rune]bool{
	'@': true, '?': true, ',': true, ':': true, '=': true, '/': true, '[': true,
	']': true, '(': true, ')': true, '{': true, '}': true, '<': true, '>': true,
}

const EOF = textscanner.EOF
