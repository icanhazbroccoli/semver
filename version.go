package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

// taken from github.com/Masterminds/semver
const SemVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

var versionRegex *regexp.Regexp

func init() {
	versionRegex = regexp.MustCompile("^" + SemVerRegex + "$")
}

// base encodes 3 10-bit numbers of a SemVer version.
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
// pre contains the pre-release tag as a string and therefore has no upper
// limitations.
type Version struct {
	base uint32
	pre  string
}

func skipTrailing(s string, i int) int {
	j := i
	for j < len(s) {
		if s[j] == ' ' || s[j] == 'v' {
			j++
			continue
		}
		break
	}
	return j
}

func readNum(s string, i int) (int, int) {
	j := i
	for j < len(s) && isNum(s[j]) {
		j++
	}
	num, err := strconv.Atoi(s[i:j])
	if err != nil {
		return -1, -1
	}
	return num, j
}

func isNum(r byte) bool {
	return r >= '0' && r <= '9'
}

func isAlpha(r byte) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDot(r byte) bool {
	return r == '.'
}

func isDash(r byte) bool {
	return r == '-'
}

func readStr(s string, i int) (string, int) {
	j := i
	for j < len(s) && (isAlpha(s[j]) || isNum(s[j]) || isDot(s[j])) {
		j++
	}
	return s[i:j], j
}

func NewVersion(s string) (*Version, error) {
	var ds [3]int
	var d int
	var pre string
	var base uint32
	var i int
	dix := 0
	i = skipTrailing(s, i)
	maxi := len(s)
	for i < maxi {
		if isNum(s[i]) {
			if dix >= len(ds) {
				goto Err
			}
			d, i = readNum(s, i)
			if i == -1 {
				goto Err
			}
			ds[dix] = d
			dix++
			if i < maxi && isDot(s[i]) {
				i++
				continue
			}
			if i < maxi && isDash(s[i]) {
				i++
				pre, _ = readStr(s, i)
				break
			}
		} else {
			goto Err
		}
	}

	for j := 0; j < 3; j++ {
		base |= (uint32(ds[j])) << uint(10*(2-j))
	}

	return &Version{
		base: base,
		pre:  pre,
	}, nil

Err:
	return nil, fmt.Errorf("failed to parse version: %q", s)
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
