package nullable

import (
	"database/sql"
	_ "unsafe"
)

type NullIgnoreScanner struct{ dst interface{} }

func NewNullIgnoreScanner(dst interface{}) *NullIgnoreScanner {
	return &NullIgnoreScanner{dst: dst}
}

func (s *NullIgnoreScanner) Scan(src interface{}) error {
	if s, ok := s.dst.(sql.Scanner); ok {
		return s.Scan(src)
	}
	if src == nil {
		return nil
	}
	return convertAssign(s.dst, src)
}

//go:linkname convertAssign database/sql.convertAssign
func convertAssign(dst, src interface{}) error
