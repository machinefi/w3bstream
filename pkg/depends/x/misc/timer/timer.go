package timer

import "time"

func Start() func() time.Duration {
	t := time.Now()
	return func() time.Duration { return time.Since(t) }
}

type Span struct {
	t    time.Time
	cost func() time.Duration
}

func StartSpan() *Span { return &Span{t: time.Now()} }

func (s *Span) Cost() time.Duration  { return time.Since(s.t) }
func (s *Span) StartedAt() time.Time { return s.t }
