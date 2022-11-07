package builder

import "context"

type comment struct {
	AdditionType
	text []byte
}

func Comment(text string) *comment { return &comment{AdditionComment, []byte(text)} }

func (c *comment) IsNil() bool { return c == nil || len(c.text) == 0 }

func (c *comment) Ex(ctx context.Context) *Ex {
	e := ExactlyExpr("")
	e.WriteComments(c.text)
	return e.Ex(ctx)
}
