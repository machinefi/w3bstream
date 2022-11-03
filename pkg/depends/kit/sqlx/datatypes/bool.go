package datatypes

import "encoding/json"

// openapi:type boolean
type Bool int

const (
	_     Bool = iota
	TRUE       // true
	FALSE      // false
)

var _ interface {
	json.Unmarshaler
	json.Marshaler
} = (*Bool)(nil)

func (v Bool) MarshalText() ([]byte, error) {
	switch v {
	case FALSE:
		return []byte("false"), nil
	case TRUE:
		return []byte("true"), nil
	default:
		return []byte("null"), nil
	}
}

func (v *Bool) UnmarshalText(data []byte) (err error) {
	switch string(data) {
	case "false":
		*v = FALSE
	case "true":
		*v = TRUE
	}
	return
}

func (v Bool) MarshalJSON() ([]byte, error) {
	return v.MarshalText()
}

func (v *Bool) UnmarshalJSON(data []byte) (err error) {
	return v.UnmarshalText(data)
}
