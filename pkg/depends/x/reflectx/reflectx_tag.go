package reflectx

import (
	"strconv"
	"strings"
)

func TagValueAndFlags(tag string) (string, map[string]bool) {
	values := strings.Split(tag, ",")
	flags := make(map[string]bool)
	if len(values[0]) > 1 {
		for _, flag := range values[1:] {
			flags[flag] = true
		}
	}
	return values[0], flags
}

type StructTag string

func (t StructTag) Name() string {
	s := string(t)
	if i := strings.Index(s, ","); i >= 0 {
		if i == 0 {
			return ""
		}
		return s[0:i]
	}
	return s
}

func (t StructTag) HasFlag(flg string) bool {
	idx := strings.Index(string(t), flg)
	return idx > 0
}

func ParseStructTag(tag string) map[string]StructTag {
	flags := map[string]StructTag{}

	for i := 0; tag != ""; {
		// skip spaces
		i = 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// meet flag name
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := tag[:i]
		tag = tag[i+1:]

		// meet flag value and unquote it
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		quoted := tag[:i+1]
		tag = tag[i+1:]
		value, err := strconv.Unquote(quoted)
		if err != nil {
			break
		}
		flags[name] = StructTag(value)
	}
	return flags
}
