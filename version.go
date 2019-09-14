package semver

// Version represents a parsed SemVer term.
// `base` encodes 3 10-bit numbers of a SemVer version.
// In binary format it looks like:
// 0b00 XXXXXXXXXX XXXXXXXXXX XXXXXXXXXX
//  |        |          |          |
//  |        |          |          '---Patch version, 10 bits
//  |        |          '---Minor version, 10 bits
//  |        '---Major version, 10 bits
//  '---Unused, 2 bits
//
// Therefore a max value for every number is 1023: the last version can't exceed
// 1023.1023.1023
//
// `pre` contains the pre-release tag as a string and therefore has no upper
// limitations.
type Version struct {
	base uint32
	pre  string
}

func NewVersion(s string) (*Version, error) {
	base, pre, err := parseVersion(s)
	if err != nil {
		return nil, err
	}

	return &Version{
		base: base,
		pre:  pre,
	}, nil
}

func (v Version) Major() uint32 {
	return (v.base >> 20) & 0x3FF
}

func (v Version) Minor() uint32 {
	return (v.base >> 10) & 0x3FF
}

func (v Version) Patch() uint32 {
	return v.base & 0x3FF
}

func (v Version) Pre() string {
	return v.pre
}

func (v Version) NextMajor() Version {
	return Version{
		base: ((v.Major() + 1) & 0x3FF) << 20,
	}
}

func (v Version) PreMajor() Version {
	return Version{
		base: ((v.Major() - 1) & 0x3FF) << 20,
	}
}

func (v Version) NextMinor() Version {
	return Version{
		base: (v.Major() & 0x3FF00000) | (((v.Minor() + 1) & 0x3FF) << 10),
	}
}

func (v Version) PrevMinor() Version {
	return Version{
		base: (v.base & 0x3FF00000) | (((v.Minor() - 1) & 0x3FF) << 10),
	}
}

func (v Version) NextPatch() Version {
	return Version{
		base: (v.base & 0x3FFFFC00) | ((v.Patch() + 1) & 0x3FF),
	}
}

func (v Version) PrevPatch() Version {
	return Version{
		base: (v.base & 0x3FFFFC00) | ((v.Patch() - 1) & 0x3FF),
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
