package mapx

type Comparable interface {
	string |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func Lt[T Comparable](a, b T) bool  { return a < b }
func Let[T Comparable](a, b T) bool { return a <= b }
func Gt[T Comparable](a, b T) bool  { return a > b }
func Get[T Comparable](a, b T) bool { return a >= b }
func Eq[T Comparable](a, b T) bool  { return a == b }
func Neq[T Comparable](a, b T) bool { return a != b }
