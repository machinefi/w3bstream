package mapx

import "github.com/pkg/errors"

type Set[K comparable] map[K]bool

func ToSet[K comparable](vals []K, conv func(K) K) (s Set[K], err error) {
	s = make(Set[K])
	for _, v := range vals {
		if _, ok := s[v]; ok && err == nil {
			err = ErrConflict
		}
		if conv != nil {
			s[conv(v)] = true
		} else {
			s[v] = true
		}
	}
	return
}

var ErrConflict = errors.Errorf("slice elements exsit conflict(s)")
