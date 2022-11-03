package ver

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Constraints struct {
	constraints [][]*constraint
}

func NewConstraint(c string) (*Constraints, error) {
	c = rewriteRange(c)

	ors := strings.Split(c, "||")
	or := make([][]*constraint, len(ors))
	for k, v := range ors {

		if !regexpValidConstraint.MatchString(v) {
			return nil, fmt.Errorf("improper constraint: %s", v)
		}

		cs := regexpFindConstraint.FindAllString(v, -1)
		if cs == nil {
			cs = append(cs, v)
		}
		result := make([]*constraint, len(cs))
		for i, s := range cs {
			pc, err := parseConstraint(s)
			if err != nil {
				return nil, err
			}

			result[i] = pc
		}
		or[k] = result
	}

	o := &Constraints{constraints: or}
	return o, nil
}

func (cs Constraints) Check(v *Version) bool {
	for _, o := range cs.constraints {
		joy := true
		for _, c := range o {
			if check, _ := c.check(v); !check {
				joy = false
				break
			}
		}

		if joy {
			return true
		}
	}

	return false
}

func (cs Constraints) Validate(v *Version) (bool, []error) {
	var e []error

	var prerelesase bool
	for _, o := range cs.constraints {
		joy := true
		for _, c := range o {
			if c.con.prerelease == "" && v.prerelease != "" {
				if !prerelesase {
					em := fmt.Errorf("%s is a prerelease version and the constraint is only looking for release versions", v)
					e = append(e, em)
					prerelesase = true
				}
				joy = false

			} else {

				if _, err := c.check(v); err != nil {
					e = append(e, err)
					joy = false
				}
			}
		}

		if joy {
			return true, []error{}
		}
	}

	return false, e
}

func (cs Constraints) String() string {
	buf := make([]string, len(cs.constraints))
	var tmp bytes.Buffer

	for k, v := range cs.constraints {
		tmp.Reset()
		vlen := len(v)
		for kk, c := range v {
			tmp.WriteString(c.string())

			if vlen > 1 && kk < vlen-1 {
				tmp.WriteString(" ")
			}
		}
		buf[k] = tmp.String()
	}

	return strings.Join(buf, " || ")
}

var (
	constraintOps           map[string]cfunc
	regexpConstraint        *regexp.Regexp
	regexpConstraintRange   *regexp.Regexp
	regexpFindConstraint    *regexp.Regexp
	regexpValidConstraint   *regexp.Regexp
	regexpConstraintVersion = regexp.MustCompile(`v?([0-9|x|X|\*]+)(\.[0-9|x|X|\*]+)?(\.[0-9|x|X|\*]+)?(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`)
)

func init() {
	constraintOps = map[string]cfunc{
		"":   constraintTildeOrEqual,
		"=":  constraintTildeOrEqual,
		"!=": constraintNotEqual,
		">":  constraintGreaterThan,
		"<":  constraintLessThan,
		">=": constraintGreaterThanEqual,
		"=>": constraintGreaterThanEqual,
		"<=": constraintLessThanEqual,
		"=<": constraintLessThanEqual,
		"~":  constraintTilde,
		"~>": constraintTilde,
		"^":  constraintCaret,
	}

	ops := make([]string, 0, len(constraintOps))
	for k := range constraintOps {
		ops = append(ops, regexp.QuoteMeta(k))
	}

	regexpConstraint = regexp.MustCompile(fmt.Sprintf(`^\s*(%s)\s*(%s)\s*$`, strings.Join(ops, "|"), regexpConstraintVersion))
	regexpConstraintRange = regexp.MustCompile(fmt.Sprintf(`\s*(%s)\s+-\s+(%s)\s*`, regexpConstraintVersion, regexpConstraintVersion))
	regexpFindConstraint = regexp.MustCompile(fmt.Sprintf(`(%s)\s*(%s)`, strings.Join(ops, "|"), regexpConstraintVersion))
	regexpValidConstraint = regexp.MustCompile(fmt.Sprintf(`^(\s*(%s)\s*(%s)\s*\,?)+$`, strings.Join(ops, "|"), regexpConstraintVersion))
}

type constraint struct {
	con *Version

	orig string

	origfunc string

	minorDirty bool
	dirty      bool
	patchDirty bool
}

func (c *constraint) check(v *Version) (bool, error) {
	return constraintOps[c.origfunc](v, c)
}

func (c *constraint) string() string {
	return c.origfunc + c.orig
}

type cfunc func(v *Version, c *constraint) (bool, error)

func parseConstraint(c string) (*constraint, error) {
	if len(c) > 0 {
		m := regexpConstraint.FindStringSubmatch(c)
		if m == nil {
			return nil, fmt.Errorf("improper constraint: %s", c)
		}

		cs := &constraint{
			orig:     m[2],
			origfunc: m[1],
		}

		ver := m[2]
		minorDirty := false
		patchDirty := false
		dirty := false
		if isX(m[3]) || m[3] == "" {
			ver = "0.0.0"
			dirty = true
		} else if isX(strings.TrimPrefix(m[4], ".")) || m[4] == "" {
			minorDirty = true
			dirty = true
			ver = fmt.Sprintf("%s.0.0%s", m[3], m[6])
		} else if isX(strings.TrimPrefix(m[5], ".")) || m[5] == "" {
			dirty = true
			patchDirty = true
			ver = fmt.Sprintf("%s%s.0%s", m[3], m[4], m[6])
		}

		con, err := ParseVersion(ver)
		if err != nil {

			return nil, errors.New("constraint Parser Error")
		}

		cs.con = con
		cs.minorDirty = minorDirty
		cs.patchDirty = patchDirty
		cs.dirty = dirty

		return cs, nil
	}

	con, err := ParseVersion("0.0.0")
	if err != nil {

		return nil, errors.New("constraint Parser Error")
	}

	cs := &constraint{
		con:        con,
		orig:       c,
		origfunc:   "",
		minorDirty: false,
		patchDirty: false,
		dirty:      true,
	}
	return cs, nil
}

func constraintNotEqual(v *Version, c *constraint) (bool, error) {
	if c.dirty {

		if v.Prerelease() != "" && c.con.Prerelease() == "" {
			return false, fmt.Errorf("%s is a prerelease version and the constraint is only looking for release versions", v)
		}

		if c.con.Major() != v.Major() {
			return true, nil
		}
		if c.con.Minor() != v.Minor() && !c.minorDirty {
			return true, nil
		} else if c.minorDirty {
			return false, fmt.Errorf("%s is equal to %s", v, c.orig)
		} else if c.con.Patch() != v.Patch() && !c.patchDirty {
			return true, nil
		} else if c.patchDirty {
			if v.Prerelease() != "" || c.con.Prerelease() != "" {
				eq := comparePrerelease(v.Prerelease(), c.con.Prerelease()) != 0
				if eq {
					return true, nil
				}
				return false, fmt.Errorf("%s is equal to %s", v, c.orig)
			}
			return false, fmt.Errorf("%s is equal to %s", v, c.orig)
		}
	}

	eq := v.Equal(c.con)
	if eq {
		return false, fmt.Errorf("%s is equal to %s", v, c.orig)
	}

	return true, nil
}

func constraintGreaterThan(v *Version, c *constraint) (bool, error) {

	if v.Prerelease() != "" && c.con.Prerelease() == "" {
		return false, fmt.Errorf("%s is a prerelease version and the constraint is only looking for release versions", v)
	}

	var eq bool

	if !c.dirty {
		eq = v.Compare(c.con) == 1
		if eq {
			return true, nil
		}
		return false, fmt.Errorf("%s is less than or equal to %s", v, c.orig)
	}

	if v.Major() > c.con.Major() {
		return true, nil
	} else if v.Major() < c.con.Major() {
		return false, fmt.Errorf("%s is less than or equal to %s", v, c.orig)
	} else if c.minorDirty {
		return false, fmt.Errorf("%s is less than or equal to %s", v, c.orig)
	} else if c.patchDirty {
		eq = v.Minor() > c.con.Minor()
		if eq {
			return true, nil
		}
		return false, fmt.Errorf("%s is less than or equal to %s", v, c.orig)
	}

	eq = v.Compare(c.con) == 1
	if eq {
		return true, nil
	}
	return false, fmt.Errorf("%s is less than or equal to %s", v, c.orig)
}

func constraintLessThan(v *Version, c *constraint) (bool, error) {
	if v.Prerelease() != "" && c.con.Prerelease() == "" {
		return false, fmt.Errorf("%s is a prerelease version and the constraint is only looking for release versions", v)
	}

	eq := v.Compare(c.con) < 0
	if eq {
		return true, nil
	}
	return false, fmt.Errorf("%s is greater than or equal to %s", v, c.orig)
}

func constraintGreaterThanEqual(v *Version, c *constraint) (bool, error) {

	if v.Prerelease() != "" && c.con.Prerelease() == "" {
		return false, fmt.Errorf("%s is a prerelease version and the constraint is only looking for release versions", v)
	}

	eq := v.Compare(c.con) >= 0
	if eq {
		return true, nil
	}
	return false, fmt.Errorf("%s is less than %s", v, c.orig)
}

func constraintLessThanEqual(v *Version, c *constraint) (bool, error) {
	if v.Prerelease() != "" && c.con.Prerelease() == "" {
		return false, fmt.Errorf("%s is a prerelease version and the constraint is only looking for release versions", v)
	}

	var eq bool

	if !c.dirty {
		eq = v.Compare(c.con) <= 0
		if eq {
			return true, nil
		}
		return false, fmt.Errorf("%s is greater than %s", v, c.orig)
	}

	if v.Major() > c.con.Major() {
		return false, fmt.Errorf("%s is greater than %s", v, c.orig)
	} else if v.Major() == c.con.Major() && v.Minor() > c.con.Minor() && !c.minorDirty {
		return false, fmt.Errorf("%s is greater than %s", v, c.orig)
	}

	return true, nil
}

func constraintTilde(v *Version, c *constraint) (bool, error) {
	if v.Prerelease() != "" && c.con.Prerelease() == "" {
		return false, fmt.Errorf("%s is a prerelease version and the constraint is only looking for release versions", v)
	}

	if v.LessThan(c.con) {
		return false, fmt.Errorf("%s is less than %s", v, c.orig)
	}

	if c.con.Major() == 0 && c.con.Minor() == 0 && c.con.Patch() == 0 &&
		!c.minorDirty && !c.patchDirty {
		return true, nil
	}

	if v.Major() != c.con.Major() {
		return false, fmt.Errorf("%s does not have same major version as %s", v, c.orig)
	}

	if v.Minor() != c.con.Minor() && !c.minorDirty {
		return false, fmt.Errorf("%s does not have same major and minor version as %s", v, c.orig)
	}

	return true, nil
}

func constraintTildeOrEqual(v *Version, c *constraint) (bool, error) {
	if v.Prerelease() != "" && c.con.Prerelease() == "" {
		return false, fmt.Errorf("%s is a prerelease version and the constraint is only looking for release versions", v)
	}

	if c.dirty {
		return constraintTilde(v, c)
	}

	eq := v.Equal(c.con)
	if eq {
		return true, nil
	}

	return false, fmt.Errorf("%s is not equal to %s", v, c.orig)
}

func constraintCaret(v *Version, c *constraint) (bool, error) {
	if v.Prerelease() != "" && c.con.Prerelease() == "" {
		return false, fmt.Errorf("%s is a prerelease version and the constraint is only looking for release versions", v)
	}

	if v.LessThan(c.con) {
		return false, fmt.Errorf("%s is less than %s", v, c.orig)
	}

	var eq bool

	if c.con.Major() > 0 || c.minorDirty {

		eq = v.Major() == c.con.Major()
		if eq {
			return true, nil
		}
		return false, fmt.Errorf("%s does not have same major version as %s", v, c.orig)
	}

	if c.con.Major() == 0 && v.Major() > 0 {
		return false, fmt.Errorf("%s does not have same major version as %s", v, c.orig)
	}
	if c.con.Minor() > 0 || c.patchDirty {
		eq = v.Minor() == c.con.Minor()
		if eq {
			return true, nil
		}
		return false, fmt.Errorf("%s does not have same minor version as %s. Expected minor versions to match when constraint major version is 0", v, c.orig)
	}

	eq = c.con.Patch() == v.Patch()
	if eq {
		return true, nil
	}
	return false, fmt.Errorf("%s does not equal %s. Expect version and constraint to equal when major and minor versions are 0", v, c.orig)
}

func isX(x string) bool {
	switch x {
	case "x", "*", "X":
		return true
	default:
		return false
	}
}

func rewriteRange(i string) string {
	m := regexpConstraintRange.FindAllStringSubmatch(i, -1)
	if m == nil {
		return i
	}
	o := i
	for _, v := range m {
		t := fmt.Sprintf(">= %s, <= %s", v[1], v[11])
		o = strings.Replace(o, v[0], t, 1)
	}

	return o
}
