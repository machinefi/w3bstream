package env

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

type Vars struct {
	Prefix string
	Values map[string]*Var
}

func LoadVarsFromEnviron(prefix string, envs []string) *Vars {
	vars := NewVars(prefix)
	for _, kv := range envs {
		pair := strings.SplitN(kv, "=", 2)
		if len(pair) != 2 {
			continue
		}
		if !strings.HasPrefix(pair[0], prefix) {
			continue
		}
		vars.Set(&Var{
			Name:  strings.Replace(pair[0], prefix+"__", "", 1),
			Value: pair[1],
		})
	}
	return vars
}

func NewVars(prefix string) *Vars { return &Vars{Prefix: prefix} }

func (vs *Vars) Set(v *Var) {
	if vs.Values == nil {
		vs.Values = map[string]*Var{}
	}
	vs.Values[v.Name] = v
}

func (vs *Vars) SetWithKeyValue(k, v string) {
	vs.Set(&Var{
		Name:  strings.Replace(k, vs.Prefix+"__", "", 1),
		Value: v,
	})
}

func (vs *Vars) Get(key string) *Var {
	if vs.Values == nil {
		return nil
	}
	return vs.Values[key]
}

func (vs *Vars) Bytes() []byte {
	kv := make(map[string]string)
	for _, v := range vs.Values {
		kv[v.Key(vs.Prefix)] = v.Value
	}
	return Env(kv)
}

func (vs *Vars) MaskBytes() []byte {
	kv := make(map[string]string)
	for _, v := range vs.Values {
		if v.Mask != "" {
			kv[v.Key(vs.Prefix)] = v.Mask
			continue
		}
		kv[v.Key(vs.Prefix)] = v.Value
	}
	return Env(kv)
}

func (vs *Vars) Len(key string) int {
	max := int64(-1)
	for _, v := range vs.Values {
		if !strings.HasPrefix(v.Name, key) {
			continue
		}
		names := strings.TrimLeft(v.Name, key+"_")
		parts := strings.Split(names, "_")
		if i, e := strconv.ParseInt(parts[0], 10, 64); e == nil && i > max {
			max = i
		}
	}
	return int(max + 1)
}

func Env(kv map[string]string) []byte {
	buf := bytes.NewBuffer(nil)

	ks := make([]string, 0)
	for k := range kv {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	for _, k := range ks {
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(kv[k])
		buf.WriteRune('\n')
	}
	return buf.Bytes()
}
