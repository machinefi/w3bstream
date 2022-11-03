package ver

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidSemVer     = errors.New("invalid semantic version")
	ErrInvalidMetadata   = errors.New("invalid Metadata string")
	ErrInvalidPrerelease = errors.New("invalid Prerelease string")
)

type Version struct {
	major         uint64
	minor         uint64
	patch         uint64
	prerelease    string
	buildMetadata string
}

var (
	regexpVerCore          = regexp.MustCompile(`(?P<major>0|[1-9]\d*)(\.(?P<minor>0|[1-9]\d*))?(\.(?P<patch>0|[1-9]\d*))?`)
	regexpVerPrerelease    = regexp.MustCompile(`(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)`)
	regexpVerBuildMetadata = regexp.MustCompile(`(?P<buildMetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*)`)
	regexpVersion          = regexp.MustCompile(`^v?` + regexpVerCore.String() + `(?:-` + regexpVerPrerelease.String() + `)?(?:\+` + regexpVerBuildMetadata.String() + `)?$`)
)

func ParseVersion(v string) (*Version, error) {
	if !regexpVersion.MatchString(v) {
		return nil, ErrInvalidSemVer
	}

	matched := regexpVersion.FindAllStringSubmatch(v, -1)[0]

	ver := Version{}

	for i, name := range regexpVersion.SubexpNames() {
		v := matched[i]

		switch name {
		case "major":
			ver.major, _ = strconv.ParseUint(v, 10, 64)
		case "minor":
			ver.minor, _ = strconv.ParseUint(v, 10, 64)
		case "patch":
			ver.patch, _ = strconv.ParseUint(v, 10, 64)
		case "prerelease":
			ver.prerelease = v
		case "buildMetadata":
			ver.buildMetadata = v
		}
	}

	return &ver, nil
}

func MustParseVersion(v string) *Version {
	sv, err := ParseVersion(v)
	if err != nil {
		panic(err)
	}
	return sv
}

func (v Version) String() string {
	buf := bytes.NewBuffer(nil)

	fmt.Fprintf(buf, "%d.%d.%d", v.major, v.minor, v.patch)

	if v.prerelease != "" {
		fmt.Fprintf(buf, "-%s", v.prerelease)
	}
	if v.buildMetadata != "" {
		fmt.Fprintf(buf, "+%s", v.buildMetadata)
	}

	return buf.String()
}

func (v Version) Major() uint64 {
	return v.major
}

func (v Version) Minor() uint64 {
	return v.minor
}

func (v Version) Patch() uint64 {
	return v.patch
}

func (v Version) Prerelease() string {
	return v.prerelease
}

func (v Version) Metadata() string {
	return v.buildMetadata
}

func (v Version) IncrPatch() *Version {
	if v.prerelease != "" {
		v.buildMetadata = ""
		v.prerelease = ""
	} else {
		v.buildMetadata = ""
		v.prerelease = ""
		v.patch = v.patch + 1
	}
	return &v
}

func (v Version) IncrMinor() *Version {
	v.buildMetadata = ""
	v.prerelease = ""
	v.patch = 0
	v.minor = v.minor + 1
	return &v
}

func (v Version) IncrMajor() *Version {
	v.buildMetadata = ""
	v.prerelease = ""
	v.patch = 0
	v.minor = 0
	v.major = v.major + 1
	return &v
}

func (v Version) WithPrerelease(prerelease string) (*Version, error) {
	if len(prerelease) > 0 {
		if !regexpVerPrerelease.MatchString(prerelease) {
			return nil, ErrInvalidPrerelease
		}
		v.prerelease = prerelease
	}
	return &v, nil
}

func (v Version) WithBuildMetadata(buildMetadata string) (*Version, error) {
	if len(buildMetadata) > 0 {
		if !regexpVerBuildMetadata.MatchString(buildMetadata) {
			return nil, ErrInvalidMetadata
		}
		v.buildMetadata = buildMetadata
	}
	return &v, nil
}

func (v *Version) LessThan(o *Version) bool {
	return v.Compare(o) < 0
}

func (v *Version) GreaterThan(o *Version) bool {
	return v.Compare(o) > 0
}

func (v *Version) Equal(o *Version) bool {
	return v.Compare(o) == 0
}

func (v *Version) Compare(o *Version) int {
	if d := compareSegment(v.Major(), o.Major()); d != 0 {
		return d
	}
	if d := compareSegment(v.Minor(), o.Minor()); d != 0 {
		return d
	}
	if d := compareSegment(v.Patch(), o.Patch()); d != 0 {
		return d
	}

	ps := v.prerelease
	po := o.Prerelease()

	if ps == "" && po == "" {
		return 0
	}
	if ps == "" {
		return 1
	}
	if po == "" {
		return -1
	}

	return comparePrerelease(ps, po)
}

func (v *Version) UnmarshalText(b []byte) error {
	ver, err := ParseVersion(string(b))
	if err != nil {
		return err
	}
	*v = *ver
	return nil
}

func (v Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

func (v *Version) UnmarshalJSON(b []byte) error { return v.UnmarshalText(b) }

func (v Version) MarshalJSON() ([]byte, error) { return v.MarshalText() }

func (v *Version) DataType(driver string) string {
	return "varchar"
}

func (v *Version) Scan(value interface{}) error {
	var s string
	s, _ = value.(string)
	ver, err := ParseVersion(s)
	if err != nil {
		return err
	}
	*v = *ver
	return nil
}

func (v Version) Value() (driver.Value, error) {
	return v.String(), nil
}

func compareSegment(v, o uint64) int {
	if v < o {
		return -1
	}
	if v > o {
		return 1
	}
	return 0
}

func comparePrerelease(v, o string) int {
	sparts := strings.Split(v, ".")
	oparts := strings.Split(o, ".")

	slen := len(sparts)
	olen := len(oparts)

	l := slen
	if olen > slen {
		l = olen
	}

	for i := 0; i < l; i++ {
		tmpStr := ""
		if i < slen {
			tmpStr = sparts[i]
		}

		otemp := ""
		if i < olen {
			otemp = oparts[i]
		}

		d := comparePrePart(tmpStr, otemp)
		if d != 0 {
			return d
		}
	}

	return 0
}

func comparePrePart(s, o string) int {
	if s == o {
		return 0
	}

	if s == "" {
		if o != "" {
			return -1
		}
		return 1
	}

	if o == "" {
		if s != "" {
			return 1
		}
		return -1
	}

	oi, n1 := strconv.ParseUint(o, 10, 64)
	si, n2 := strconv.ParseUint(s, 10, 64)

	if n1 != nil && n2 != nil {
		if s > o {
			return 1
		}
		return -1
	} else if n1 != nil {
		return -1
	} else if n2 != nil {
		return 1
	}
	if si > oi {
		return 1
	}
	return -1

}
