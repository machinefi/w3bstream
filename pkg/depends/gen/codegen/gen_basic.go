package codegen

import (
	"bytes"
	"fmt"
	"go/token"
	"sort"
	"strconv"
	"strings"

	"github.com/machinefi/w3bstream/pkg/depends/x/stringsx"
)

type (
	Snippet interface{ Bytes() []byte }

	SnippetSpec interface {
		Snippet
		IfSpec
	}

	SnippetType interface {
		Snippet
		IfType
	}

	SnippetCanAddr interface {
		Snippet
		IfCanAddr
	}

	IfSpec          interface{ _spec() }
	IfCanAddr       interface{ _canAddr() }
	IfCanBeIfMethod interface{ _canBeInterfaceMethod() }
	IfType          interface{ _type() }

	FnAlaise func(string) string
)

type SnippetLiteral string

var _ Snippet = SnippetLiteral("")

func (s SnippetLiteral) Bytes() []byte { return []byte(s) }

type SnippetLiteralCompose struct {
	Type     SnippetType
	Elements []Snippet
}

var _ Snippet = (*SnippetLiteralCompose)(nil)

func (s *SnippetLiteralCompose) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	if s.Type != nil {
		buf.Write(s.Type.Bytes())
	}

	buf.WriteRune('{')

	for _, e := range s.Elements {
		buf.WriteRune('\n')
		buf.Write(e.Bytes())
		buf.WriteRune(',')
	}

	buf.WriteRune('\n')
	buf.WriteRune('}')
	return buf.Bytes()
}

/*
SnippetBlock code block, like
```go

	var a = "Hello CodeGen"
	fmt.Print(a)

```
*/
type SnippetBlock []Snippet

var _ Snippet = SnippetBlock(nil)

func (s SnippetBlock) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	for _, sn := range s {
		if sn == nil {
			continue
		}
		buf.Write(sn.Bytes())
		buf.WriteRune('\n')
	}

	return buf.Bytes()
}

/*
SnippetBlockWithBrace code block quote with '{' and '}', like
```go

	{
		var a = "Hello CodeGen"
		fmt.Print(a)
	}

```
*/
type SnippetBlockWithBrace []Snippet

var _ Snippet = SnippetBlockWithBrace(nil)

func (s SnippetBlockWithBrace) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteRune('{')
	buf.WriteRune('\n')
	for _, sn := range s {
		if sn == nil {
			continue
		}
		buf.Write(sn.Bytes())
		buf.WriteRune('\n')
	}
	buf.WriteRune('}')

	return buf.Bytes()
}

// SnippetBuiltIn built-in symbols, like `int`, `println`
type SnippetBuiltIn string

var _ Snippet = SnippetBuiltIn("")

func (s SnippetBuiltIn) Bytes() []byte { return []byte(s) }

// SnippetIdent idenifier `a := 0 // a is a identifier`
type SnippetIdent string

var _ SnippetCanAddr = SnippetIdent("")

func (s SnippetIdent) Bytes() []byte { return []byte(s) }

func (s SnippetIdent) _canAddr() {}

func (s SnippetIdent) UpperCamelCase() *SnippetIdent {
	return Ident(stringsx.UpperCamelCase(string(s)))
}

func (s SnippetIdent) LowerCamelCase() *SnippetIdent {
	return Ident(stringsx.LowerCamelCase(string(s)))
}

func (s SnippetIdent) UpperSnakeCase() *SnippetIdent {
	return Ident(stringsx.UpperSnakeCase(string(s)))
}

func (s SnippetIdent) LowerSnakeCase() *SnippetIdent {
	return Ident(stringsx.LowerSnakeCase(string(s)))
}

// SnippetComments comment code
type SnippetComments struct {
	OneLine  bool
	Comments []string
}

var _ Snippet = (*SnippetComments)(nil)

func (s *SnippetComments) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	if s.IsOneLine() {
		buf.WriteString("// ")
		buf.WriteString(s.Comments[0])
	} else {
		for i, c := range s.Comments {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString("// ")
			buf.WriteString(c)
		}
	}
	return buf.Bytes()
}

func (s SnippetComments) AsOneLine() *SnippetComments { s.OneLine = true; return &s }

func (s SnippetComments) IsOneLine() bool { return s.OneLine && len(s.Comments) == 1 }

func (s SnippetComments) Append(cmt ...string) SnippetComments {
	for _, c := range cmt {
		s.Comments = append(s.Comments, strings.Split(c, "\n")...)
	}
	return s
}

// SnippetExpr expression `a == 0 // a == 0 is a expression`
type SnippetExpr string

var _ Snippet = SnippetExpr("")

func (s SnippetExpr) Bytes() []byte { return []byte(s) }

type SnippetKVExpr struct {
	Snippet
	K, V Snippet
}

var _ Snippet = (*SnippetKVExpr)(nil)

func (s *SnippetKVExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.Write(s.K.Bytes())
	buf.WriteString(token.COLON.String())
	buf.WriteRune(' ')
	buf.Write(s.V.Bytes())

	return buf.Bytes()
}

type SnippetTypeDecl struct {
	Token token.Token
	Specs []SnippetSpec
	*SnippetComments
}

var _ Snippet = (*SnippetTypeDecl)(nil)

func (s *SnippetTypeDecl) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	mul := len(s.Specs) > 1

	if s.SnippetComments != nil {
		buf.Write(s.SnippetComments.Bytes())
		buf.WriteRune('\n')
	}

	buf.WriteString(s.Token.String())
	buf.WriteRune(' ')

	if mul {
		buf.WriteRune('(')
		buf.WriteRune('\n')
	}

	for i, spec := range s.Specs {
		if i > 0 {
			buf.WriteRune('\n')
		}
		buf.Write(spec.Bytes())
	}

	if mul {
		buf.WriteRune('\n')
		buf.WriteRune(')')
	}

	return buf.Bytes()
}

func (s SnippetTypeDecl) WithComments(comments ...string) *SnippetTypeDecl {
	if len(comments) > 0 {
		s.SnippetComments = Comments(comments...)
	}
	return &s
}

// SnippetField define a field or var
// eg:
// a int
// a, b int
// AliasString = string
type SnippetField struct {
	SnippetSpec
	SnippetCanAddr

	Type  SnippetType
	Names []*SnippetIdent
	Tag   string
	Alias bool

	*SnippetComments
}

var _ Snippet = (*SnippetField)(nil)

func (s *SnippetField) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	if s.SnippetComments != nil && !s.SnippetComments.IsOneLine() {
		tmp := s.SnippetComments.Bytes()
		buf.Write(tmp)
		if len(tmp) > 0 {
			buf.WriteRune('\n')
		}
	}

	for i := range s.Names {
		if i > 0 {
			buf.WriteRune(',')
			buf.WriteRune(' ')
		}
		buf.Write(s.Names[i].Bytes())
	}

	if len(s.Names) > 0 {
		if s.Alias {
			buf.WriteRune(' ')
			buf.WriteRune('=')
		}
	}

	// type inference
	if s.Type != nil {
		if len(s.Names) > 0 {
			buf.WriteRune(' ')
		}
		buf.Write(s.Type.Bytes())
	}

	if s.Tag != "" {
		buf.WriteRune(' ')
		buf.WriteRune('`')
		buf.WriteString(s.Tag)
		buf.WriteRune('`')
	}

	if s.SnippetComments != nil && s.SnippetComments.IsOneLine() {
		buf.WriteRune(' ')
		buf.Write(s.SnippetComments.Bytes())
	}

	return buf.Bytes()
}

func (s SnippetField) WithTag(tag string) *SnippetField {
	s.Tag = tag
	return &s
}

func (s SnippetField) WithTags(tags map[string][]string) *SnippetField {
	buf := bytes.NewBuffer(nil)

	names := make([]string, 0)
	for tag := range tags {
		names = append(names, tag)
	}
	sort.Strings(names)

	for i, tag := range names {
		if i > 0 {
			buf.WriteRune(' ')
		}
		values := make([]string, 0)
		for j := range tags[tag] {
			v := tags[tag][j]
			if v != "" {
				values = append(values, v)
			}
		}
		buf.WriteString(tag)
		buf.WriteRune(':')
		buf.WriteString(strconv.Quote(strings.Join(values, ",")))
	}

	s.Tag = buf.String()
	return &s
}

func (s SnippetField) WithoutTag() *SnippetField {
	s.Tag = ""
	return &s
}

func (s SnippetField) WithComments(cmt ...string) *SnippetField {
	s.SnippetComments = Comments(cmt...)
	return &s
}

func (s SnippetField) WithOneLineComment(cmt string) *SnippetField {
	if s.SnippetComments == nil {
		s.SnippetComments = &SnippetComments{}
	}
	s.SnippetComments.OneLine = true
	s.SnippetComments.Comments = []string{cmt}
	return &s
}

func (s SnippetField) AsAlias() *SnippetField {
	s.Alias = true
	return &s
}

type SnippetCaseClause struct {
	Case []Snippet
	Blk  SnippetBlock
}

var _ Snippet = (*SnippetCaseClause)(nil)

func (s *SnippetCaseClause) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	if len(s.Case) == 0 {
		buf.WriteString(token.DEFAULT.String())
	} else {
		buf.WriteString(token.CASE.String() + " ")
		for i, c := range s.Case {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.Write(c.Bytes())
		}
	}

	buf.WriteRune(':')
	for _, b := range s.Blk {
		buf.WriteRune('\n')
		buf.Write(b.Bytes())
	}

	return buf.Bytes()
}

func (s SnippetCaseClause) Do(bodies ...Snippet) *SnippetCaseClause {
	s.Blk = bodies
	return &s
}

type SnippetTypeAssertExpr struct {
	Asserter Snippet
	Type     SnippetType
}

var _ Snippet = (*SnippetTypeAssertExpr)(nil)

func (s *SnippetTypeAssertExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.Write(s.Asserter.Bytes())

	buf.WriteRune('.')
	buf.WriteRune('(')
	buf.Write(s.Type.Bytes())
	buf.WriteRune(')')

	return buf.Bytes()
}

type SnippetAssignStmt struct {
	IfSpec
	Token token.Token
	Ls    []SnippetCanAddr
	Rs    []Snippet
}

var _ Snippet = (*SnippetAssignStmt)(nil)

func (s *SnippetAssignStmt) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	for i, l := range s.Ls {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.Write(l.Bytes())
	}

	if len(s.Rs) > 0 {
		buf.WriteRune(' ')
		buf.WriteString(s.Token.String())
		buf.WriteRune(' ')

		for i, r := range s.Rs {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.Write(r.Bytes())
		}
	}

	return buf.Bytes()
}

func (s SnippetAssignStmt) By(rs ...Snippet) *SnippetAssignStmt {
	s.Rs = rs
	return &s
}

type SnippetReturnStmt struct {
	Res []Snippet
}

var _ Snippet = (*SnippetReturnStmt)(nil)

func (s *SnippetReturnStmt) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("return")

	for i, r := range s.Res {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune(' ')
		buf.Write(r.Bytes())
	}

	return buf.Bytes()
}

type SnippetSelectStmt struct {
	Clauses []*SnippetCaseClause
}

var _ Snippet = (*SnippetSelectStmt)(nil)

func (s *SnippetSelectStmt) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.SELECT.String())
	buf.WriteRune(' ')
	buf.WriteRune('{')

	for _, c := range s.Clauses {
		buf.WriteRune('\n')
		buf.Write(c.Bytes())
	}

	buf.WriteRune('\n')
	buf.WriteRune('}')

	return buf.Bytes()
}

type SnippetSwitchStmt struct {
	Init, Cond Snippet
	Clauses    []*SnippetCaseClause
}

var _ Snippet = (*SnippetSwitchStmt)(nil)

func (s *SnippetSwitchStmt) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.SWITCH.String())
	if s.Cond != nil {
		if s.Init != nil {
			buf.WriteRune(' ')
			buf.Write(s.Init.Bytes())
			buf.WriteRune(';')
		}
		buf.WriteRune(' ')
		buf.Write(s.Cond.Bytes())
	}

	buf.WriteRune(' ')
	buf.WriteRune('{')

	for _, c := range s.Clauses {
		buf.WriteRune('\n')
		buf.Write(c.Bytes())
	}

	buf.WriteRune('\n')
	buf.WriteRune('}')

	return buf.Bytes()
}

func (s SnippetSwitchStmt) InitWith(init Snippet) *SnippetSwitchStmt {
	s.Init = init
	return &s
}

func (s SnippetSwitchStmt) When(clauses ...*SnippetCaseClause) *SnippetSwitchStmt {
	s.Clauses = append(s.Clauses, clauses...)
	return &s
}

type SnippetForRangeStmt struct {
	K, V   SnippetIdent
	Ranger Snippet
	Blk    SnippetBlockWithBrace
}

var _ Snippet = (*SnippetForRangeStmt)(nil)

func (s *SnippetForRangeStmt) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.FOR.String())
	buf.WriteRune(' ')

	if s.K != AnonymousIdent || s.V != AnonymousIdent {
		if s.K != AnonymousIdent && s.V == AnonymousIdent {
			buf.Write(s.K.Bytes())
		} else if s.V != AnonymousIdent {
			buf.Write(s.K.Bytes())
			buf.WriteRune(',')
			buf.WriteRune(' ')
			buf.Write(s.V.Bytes())
		}
		buf.WriteRune(' ')
		buf.WriteString(token.DEFINE.String())
		buf.WriteRune(' ')
	}

	buf.WriteString(token.RANGE.String())
	buf.WriteRune(' ')
	buf.Write(s.Ranger.Bytes())
	buf.WriteRune(' ')
	buf.Write(s.Blk.Bytes())

	return buf.Bytes()
}

func (s SnippetForRangeStmt) Do(blk ...Snippet) *SnippetForRangeStmt {
	s.Blk = blk
	return &s
}

type SnippetForStmt struct {
	Init, Cond, Post Snippet
	Blk              SnippetBlockWithBrace
}

var _ Snippet = (*SnippetForStmt)(nil)

func (s *SnippetForStmt) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.FOR.String())
	if s.Init != nil {
		buf.WriteRune(' ')
		buf.Write(s.Init.Bytes())
		buf.WriteRune(';')
	}
	if s.Cond != nil {
		buf.WriteRune(' ')
		buf.Write(s.Cond.Bytes())
	}
	if s.Post != nil {
		buf.WriteRune(';')
		buf.WriteRune(' ')
		buf.Write(s.Post.Bytes())
	}
	buf.WriteRune(' ')
	buf.Write(s.Blk.Bytes())

	return buf.Bytes()
}

func (s SnippetForStmt) Do(blk ...Snippet) *SnippetForStmt {
	s.Blk = blk
	return &s
}

type SnippetIfStmt struct {
	Init, Cond Snippet
	Blk        SnippetBlockWithBrace
	ElseList   []*SnippetIfStmt
}

var _ Snippet = (*SnippetIfStmt)(nil)

func (s *SnippetIfStmt) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	if s.Cond != nil {
		buf.WriteString(token.IF.String())
	}

	if s.Init != nil {
		buf.WriteRune(' ')
		buf.Write(s.Init.Bytes())
		buf.WriteRune(';')
	}

	if s.Cond != nil {
		buf.WriteRune(' ')
		buf.Write(s.Cond.Bytes())
	}

	buf.WriteRune(' ')
	buf.Write(s.Blk.Bytes())

	for _, then := range s.ElseList {
		buf.WriteRune(' ')
		buf.WriteString(token.ELSE.String())
		if then.Cond != nil {
			buf.WriteRune(' ')
		}
		buf.Write(then.Bytes())
	}

	return buf.Bytes()
}

func (s SnippetIfStmt) InitWith(init Snippet) *SnippetIfStmt {
	s.Init = init
	return &s
}

func (s SnippetIfStmt) Do(ss ...Snippet) *SnippetIfStmt {
	s.Blk = append(s.Blk, ss...)
	return &s
}

func (s SnippetIfStmt) Else(sub *SnippetIfStmt) *SnippetIfStmt {
	s.ElseList = append(s.ElseList, sub)
	return &s
}

type SnippetRefExpr struct {
	Lead Snippet
	Refs []Snippet
}

var _ Snippet = (*SnippetRefExpr)(nil)

func (s *SnippetRefExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.Write(s.Lead.Bytes())
	for _, ref := range s.Refs {
		buf.WriteRune('.')
		buf.Write(ref.Bytes())
	}

	return buf.Bytes()
}

type SnippetStarExpr struct {
	SnippetCanAddr
	SnippetType
	T SnippetType
}

var _ Snippet = (*SnippetStarExpr)(nil)

func (s *SnippetStarExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.MUL.String())
	buf.Write(s.T.Bytes())

	return buf.Bytes()
}

type SnippetAccessValueExpr struct {
	V SnippetCanAddr
	IfCanAddr
}

func (s *SnippetAccessValueExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(token.MUL.String())
	buf.WriteRune('(')
	buf.Write(s.V.Bytes())
	buf.WriteRune(')')
	return buf.Bytes()
}

type SnippetAddrExpr struct {
	V SnippetCanAddr
}

var _ Snippet = (*SnippetAddrExpr)(nil)

func (s *SnippetAddrExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.AND.String())
	buf.Write(s.V.Bytes())

	return buf.Bytes()
}

type SnippetParenExpr struct {
	SnippetCanAddr
	V Snippet
}

var _ Snippet = (*SnippetParenExpr)(nil)

func (s *SnippetParenExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteRune('(')
	buf.Write(s.V.Bytes())
	buf.WriteRune(')')

	return buf.Bytes()
}

type SnippetIncExpr struct {
	Value SnippetCanAddr
}

var _ Snippet = (*SnippetIncExpr)(nil)

func (s *SnippetIncExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.Write(s.Value.Bytes())
	buf.WriteString(token.INC.String())

	return buf.Bytes()
}

type SnippetDecExpr struct {
	Value SnippetCanAddr
}

var _ Snippet = (*SnippetDecExpr)(nil)

func (s *SnippetDecExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.Write(s.Value.Bytes())
	buf.WriteString(token.DEC.String())

	return buf.Bytes()
}

type SnippetArrowExpr struct {
	SnippetCanAddr
	Chan Snippet
}

var _ Snippet = (*SnippetArrowExpr)(nil)

func (s *SnippetArrowExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.ARROW.String())
	buf.Write(s.Chan.Bytes())

	return buf.Bytes()
}

type SnippetAccessExpr struct {
	V     Snippet
	Index Snippet
}

var _ Snippet = (*SnippetAccessExpr)(nil)

func (s *SnippetAccessExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.Write(s.V.Bytes())
	buf.WriteString(token.LBRACK.String())
	buf.Write(s.Index.Bytes())
	buf.WriteString(token.RBRACK.String())

	return buf.Bytes()
}

type SnippetCallExpr struct {
	Callee   Snippet
	Args     []Snippet
	Ellipsis bool
	Modifier token.Token
}

var _ Snippet = (*SnippetCallExpr)(nil)

func (s *SnippetCallExpr) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	if s.Modifier > token.ILLEGAL {
		buf.WriteString(s.Modifier.String())
		buf.WriteRune(' ')
	}

	buf.Write(s.Callee.Bytes())

	buf.WriteRune('(')
	for i, a := range s.Args {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.Write(a.Bytes())
	}

	if s.Ellipsis {
		buf.WriteString(token.ELLIPSIS.String())
	}
	buf.WriteRune(')')

	return buf.Bytes()
}

func (s SnippetCallExpr) AsDefer() *SnippetCallExpr {
	s.Modifier = token.DEFER
	return &s
}

func (s SnippetCallExpr) AsRoutine() *SnippetCallExpr {
	s.Modifier = token.GO
	return &s
}

func (s SnippetCallExpr) WithEllipsis() *SnippetCallExpr {
	s.Ellipsis = true
	return &s
}

type BuiltInType string

var _ SnippetType = BuiltInType("")

func (t BuiltInType) Bytes() []byte { return []byte(t) }

func (t BuiltInType) _type() {}

type MapType struct {
	SnippetType
	Tk, Tv SnippetType
}

var _ SnippetType = (*MapType)(nil)

func (m *MapType) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.MAP.String())
	buf.WriteRune('[')
	buf.Write(m.Tk.Bytes())
	buf.WriteRune(']')
	buf.Write(m.Tv.Bytes())

	return buf.Bytes()
}

type ArrayType struct {
	SnippetType
	T   SnippetType
	Len int
}

var _ SnippetType = (*ArrayType)(nil)

func (s *ArrayType) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(fmt.Sprintf("[%d]", s.Len))
	buf.Write(s.T.Bytes())

	return buf.Bytes()
}

type SliceType struct {
	SnippetType
	T SnippetType
}

var _ SnippetType = (*SliceType)(nil)

func (s *SliceType) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteRune('[')
	buf.WriteRune(']')
	buf.Write(s.T.Bytes())

	return buf.Bytes()
}

type ChanMode uint8

const (
	ChanModeRO ChanMode = 0x01
	ChanModeWO ChanMode = 0x10
	ChanModeRW ChanMode = 0x11
)

type ChanType struct {
	SnippetType
	T    SnippetType
	Mode ChanMode
}

var _ SnippetType = (*ChanType)(nil)

func (s *ChanType) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	if s.Mode == ChanModeRO {
		buf.WriteString(token.ARROW.String())
	}
	buf.WriteString(token.CHAN.String())
	if s.Mode == ChanModeWO {
		buf.WriteString(token.ARROW.String())
	}
	buf.WriteRune(' ')
	buf.Write(s.T.Bytes())

	return buf.Bytes()
}

type EllipsisType struct {
	SnippetType
	T SnippetType
}

var _ SnippetType = (*EllipsisType)(nil)

func (s *EllipsisType) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.ELLIPSIS.String())
	buf.Write(s.T.Bytes())

	return buf.Bytes()
}

type FuncType struct {
	SnippetType
	IfCanBeIfMethod
	Name    *SnippetIdent
	Recv    *SnippetField
	Args    []*SnippetField
	Rets    []*SnippetField
	Blk     SnippetBlockWithBrace
	noToken bool
}

var _ SnippetType = (*FuncType)(nil)

func (f *FuncType) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	if !f.noToken {
		buf.WriteString(token.FUNC.String())
		buf.WriteRune(' ')
	}

	if f.Recv != nil {
		buf.WriteRune('(')
		buf.Write(f.Recv.Bytes())
		buf.WriteRune(')')
		buf.WriteRune(' ')
	}

	if f.Name != nil {
		buf.Write(f.Name.Bytes())
	}

	buf.WriteByte('(')
	for i := range f.Args {
		if i > 0 {
			buf.WriteRune(',')
			buf.WriteRune(' ')
		}
		buf.Write(f.Args[i].WithoutTag().Bytes())
	}
	buf.WriteByte(')')
	buf.WriteRune(' ')

	quoteRet := len(f.Rets) > 0 && len(f.Rets[0].Names) > 0 || len(f.Rets) > 1

	if quoteRet {
		buf.WriteRune('(')
	}

	for i := range f.Rets {
		if i > 0 {
			buf.WriteRune(',')
			buf.WriteRune(' ')
		}
		buf.Write(f.Rets[i].WithoutTag().Bytes())
	}

	if quoteRet {
		buf.WriteRune(')')
		buf.WriteRune(' ')
	} else {
		if len(f.Rets) > 0 {
			buf.WriteRune(' ')
		}
	}

	if f.Blk != nil {
		buf.Write(f.Blk.Bytes())
	}
	return buf.Bytes()
}

func (f FuncType) WithoutToken() *FuncType { f.noToken = true; return &f }

func (f FuncType) Named(name string) *FuncType { f.Name = Ident(name); return &f }

func (f FuncType) MethodOf(rcv *SnippetField) *FuncType { f.Recv = rcv; return &f }

func (f FuncType) Return(rets ...*SnippetField) *FuncType { f.Rets = rets; return &f }

func (f FuncType) Do(ss ...Snippet) *FuncType { f.Blk = append([]Snippet{}, ss...); return &f }

type StructType struct {
	SnippetType
	Fields []*SnippetField
}

var _ SnippetType = (*StructType)(nil)

func (s *StructType) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.STRUCT.String())

	buf.WriteRune(' ')
	buf.WriteRune('{')

	for i := range s.Fields {
		buf.WriteRune('\n')
		buf.Write(s.Fields[i].Bytes())
	}

	buf.WriteRune('\n')
	buf.WriteRune('}')

	return buf.Bytes()
}

type InterfaceType struct {
	SnippetType
	Methods []IfCanBeIfMethod
}

var _ SnippetType = (*InterfaceType)(nil)

func (s *InterfaceType) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString(token.INTERFACE.String())
	if len(s.Methods) == 0 {
		buf.WriteRune('{')
		buf.WriteRune('}')
		return buf.Bytes()
	}
	buf.WriteRune(' ')
	buf.WriteRune('{')
	buf.WriteRune('\n')
	for i := range s.Methods {
		switch method := s.Methods[i].(type) {
		case *FuncType:
			buf.Write(method.WithoutToken().Bytes())
		case *NamedType:
			buf.Write(method.Bytes())
		}
		buf.WriteRune('\n')
	}
	buf.WriteRune('}')

	return buf.Bytes()
}

type NamedType struct {
	SnippetType
	IfCanBeIfMethod
	IfCanAddr
	Name *SnippetIdent
}

var _ SnippetType = (*NamedType)(nil)

func (t *NamedType) Bytes() []byte { return t.Name.Bytes() }
