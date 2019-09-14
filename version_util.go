package semver

import (
	"fmt"
	"strconv"
)

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

func parseVersion(s string) (uint32, string, error) {
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
	return base, pre, nil

Err:
	return 0, "", fmt.Errorf("failed to parse version: %q", s)
}
