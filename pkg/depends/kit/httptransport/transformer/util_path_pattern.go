package transformer

import (
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func ParamsFromMap(m map[string]string) httprouter.Params {
	params := httprouter.Params{}
	for k, v := range m {
		params = append(params, httprouter.Param{
			Key:   k,
			Value: v,
		})
	}
	return params
}

func NewPathnamePattern(pattern string) *PathnamePattern {
	parts := ToPathParts(pattern)

	idxKeys := map[int]string{}

	for i, p := range parts {
		if p[0] == ':' {
			idxKeys[i] = p[1:]
		}
	}

	return &PathnamePattern{
		parts,
		idxKeys,
	}
}

type PathnamePattern struct {
	parts   []string
	idxKeys map[int]string
}

func (pp *PathnamePattern) String() string { return "/" + strings.Join(pp.parts, "/") }

func (pp *PathnamePattern) Stringify(params httprouter.Params) string {
	if len(pp.idxKeys) == 0 {
		return pp.String()
	}

	parts := append([]string{}, pp.parts...)

	for idx, key := range pp.idxKeys {
		v := params.ByName(key)
		if v == "" {
			v = "-"
		}
		parts[idx] = v
	}

	return (&PathnamePattern{parts: parts}).String()
}

func (pp *PathnamePattern) Parse(pathname string) (params httprouter.Params, err error) {
	parts := ToPathParts(pathname)

	if len(parts) != len(pp.parts) {
		return nil, errors.Errorf("pathname %s is not match %s", pathname, pp)
	}

	for idx, part := range pp.parts {
		if key, ok := pp.idxKeys[idx]; ok {
			params = append(params, httprouter.Param{
				Key:   key,
				Value: parts[idx],
			})
		} else if part != parts[idx] {
			return nil, errors.Errorf("pathname %s is not match %s", pathname, pp)
		}
	}

	return
}

func ToPathParts(p string) []string {
	p = httprouter.CleanPath(p)
	if p[0] == '/' {
		p = p[1:]
	}
	if p == "" {
		return make([]string, 0)
	}
	return strings.Split(p, "/")
}
