package semver

import (
	"fmt"
)

func parseVersion(s string) (uint32, string, error) {
	var ds [3]uint32
	var d int
	var pre string
	var base uint32
	dix := 0
	i, maxi := 0, len(s)
	i = skipTrailing(s, i)
	for i < maxi {
		if isNum(s[i]) {
			if dix >= len(ds) {
				goto Err
			}
			d, i = readNum(s, i)
			if i == -1 {
				goto Err
			}
			ds[dix] = uint32(d)
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

	base = (ds[0] << 20) | (ds[1] << 10) | ds[2]
	return base, pre, nil

Err:
	return 0, "", fmt.Errorf("failed to parse version: %q", s)
}
