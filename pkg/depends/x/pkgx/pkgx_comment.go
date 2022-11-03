package pkgx

import (
	"go/ast"
	"go/token"
	"sort"
)

type CommentScanner struct {
	*ast.File
	ast.CommentMap
}

func NewCommentScanner(fs *token.FileSet, f *ast.File) *CommentScanner {
	return &CommentScanner{
		File:       f,
		CommentMap: ast.NewCommentMap(fs, f, f.Comments),
	}
}

func (c *CommentScanner) CommentsOf(n ast.Node) string {
	return StringifyCommentGroup(c.CommentGroupsOf(n)...)
}

func (c *CommentScanner) CommentGroupsOf(n ast.Node) (cgs []*ast.CommentGroup) {
	if n == nil {
		return
	}

	switch n.(type) {
	case *ast.File, *ast.Field, ast.Stmt, ast.Decl:
		if comments, ok := c.CommentMap[n]; ok {
			cgs = comments
		}
	case ast.Spec:
		if comments, ok := c.CommentMap[n]; ok {
			cgs = append(cgs, comments...)
		}
		if len(cgs) == 0 {
			for node, comments := range c.CommentMap {
				if decl, ok := node.(*ast.GenDecl); ok {
					for _, spec := range decl.Specs {
						if n == spec {
							cgs = append(cgs, comments...)
						}
					}
				}
			}
		}
	default:
		var (
			pos    token.Pos = -1
			parent ast.Node
		)
		ast.Inspect(c.File, func(node ast.Node) bool {
			switch node.(type) {
			case *ast.Field, ast.Stmt, ast.Decl, ast.Spec:
				if n.Pos() >= node.Pos() && n.End() <= node.End() {
					next := n.Pos() - node.Pos()
					if pos == -1 || next <= pos {
						pos, parent = next, node
					}
				}
			}
			return true
		})
		if parent != nil {
			cgs = c.CommentGroupsOf(parent)
		}
	}
	sort.Slice(cgs, func(i, j int) bool {
		return cgs[i].Pos() < cgs[j].Pos()
	})
	return
}
