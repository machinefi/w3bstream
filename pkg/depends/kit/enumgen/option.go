package enumgen

import (
	"sort"
	"strconv"
)

type Option struct {
	Label string   `json:"label"`
	Str   *string  `json:"str,omitempty"`
	Int   *int64   `json:"int,omitempty"`
	Float *float64 `json:"float,omitempty"`
}

func (o Option) Value() interface{} {
	if o.Str != nil {
		return *o.Str
	}
	if o.Int != nil {
		return *o.Int
	}
	if o.Float != nil {
		return *o.Float
	}
	return nil
}

type Options []Option

var _ sort.Interface = Options{}

func (o Options) Len() int { return len(o) }

func (o Options) Values() []interface{} {
	values := make([]interface{}, len(o))
	for i, v := range o {
		values[i] = v.Value()
	}
	return values
}

func (o Options) Less(i, j int) bool {
	if o[i].Float != nil {
		return *o[i].Float < *o[j].Float
	}
	if o[i].Int != nil {
		return *o[i].Int < *o[j].Int
	}
	return *o[i].Str < *o[j].Str
}

func (o Options) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// NewOption int-string option
func NewOption(i int64, str, label string) *Option {
	ret := &Option{
		Label: label,
		Str:   &str,
		Int:   &i,
	}
	if label == "" {
		ret.Label = str
	}
	return ret
}

func NewIntOption(i int64, label string) *Option {
	ret := &Option{
		Label: label,
		Int:   &i,
	}
	if label == "" {
		ret.Label = strconv.FormatInt(i, 10)
	}
	return ret
}

func NewFloatOption(f float64, label string) *Option {
	ret := &Option{
		Label: label,
		Float: &f,
	}
	if label == "" {
		ret.Label = strconv.FormatFloat(f, 'f', -1, 64)
	}
	return ret
}

func NewStringOption(str, label string) *Option {
	ret := &Option{
		Label: label,
		Str:   &str,
	}
	if label == "" {
		ret.Label = str
	}
	return ret
}
