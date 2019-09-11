package semver

import (
	"regexp"
	"strconv"
	"strings"
)

// taken from github.com/Masterminds/semver
const SemVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

var versionRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile("^" + SemVerRegex + "$")
}

type Version struct {
	base uint32
	pre  string
}

func NewVersion(v string) (*Version, error) {
	m := versionRegex.FindStringSubmatch(v)
	if m == nil {
		return nil, ErrInvalidSemVer
	}
	var base uint32
	if v, err := strconv.Atoi(m[1]); err != nil {
		return nil, err
	} else {
		base |= (uint32(v) & 0xFF) << 16
	}
	if v, err := strconv.Atoi(strings.TrimPrefix(m[2], ".")); err != nil {
		return nil, err
	} else {
		base |= (uint32(v) & 0xFF) << 8
	}
	if v, err := strconv.Atoi(strings.TrimPrefix(m[3], ".")); err != nil {
		return nil, err
	} else {
		base |= uint32(v) & 0xFF
	}
	return &Version{
		base: base,
		pre:  m[5],
	}, nil
}

func (v Version) Major() uint32 {
	return (v.base >> 16) & 0xFF
}

func (v Version) Minor() uint32 {
	return (v.base >> 8) & 0xFF
}

func (v Version) Patch() uint32 {
	return v.base & 0xFF
}

func (v Version) Pre() string {
	return v.pre
}

func (v Version) NextMajor() Version {
	return Version{
		base: ((v.Major() + 1) & 0xFF) << 16,
	}
}

func (v Version) PreMajor() Version {
	return Version{
		base: ((v.Major() - 1) & 0xFF) << 16,
	}
}

func (v Version) NextMinor() Version {
	return Version{
		base: (v.Major() & 0xFF0000) | (((v.Minor() + 1) & 0xFF) << 8),
	}
}

func (v Version) PrevMinor() Version {
	return Version{
		base: (v.base & 0xFF0000) | (((v.Minor() - 1) & 0xFF) << 8),
	}
}

func (v Version) NextPatch() Version {
	return Version{
		base: (v.base & 0xFFFF00) | ((v.Patch() + 1) & 0xFF),
	}
}

func (v Version) PrevPatch() Version {
	return Version{
		base: (v.base & 0xFFFF00) | ((v.Patch() - 1) & 0xFF),
	}
}

func (v1 Version) Equal(v2 *Version) bool {
	return v1.base == v2.base && v1.pre == v2.pre
}

func (v1 Version) Less(v2 *Version) bool {
	if v1.base < v2.base {
		return true
	} else if v1.base == v2.base {
		lv1, lv2 := len(v1.pre), len(v2.pre)
		if lv1 != 0 && lv2 != 0 {
			return v1.pre < v2.pre
		}
		return lv1 > 0
	}
	return false
}
