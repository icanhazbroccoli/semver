package semver

import "strconv"

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

func isStar(r byte) bool {
	return r == '*' || r == 'x' || r == 'X'
}

func isOpChar(r byte) bool {
	return r == '=' || r == '<' || r == '>' || r == '^' || r == '!' || r == '~'
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
	j, maxj := i, len(s)
	for j < maxj && isNum(s[j]) {
		j++
	}
	num, err := strconv.Atoi(s[i:j])
	if err != nil {
		return -1, -1
	}
	return num, j
}

func readStr(s string, i int) (string, int) {
	j, maxj := i, len(s)
	for j < maxj && (isAlpha(s[j]) || isNum(s[j]) || isDot(s[j])) {
		j++
	}
	return s[i:j], j
}

func readOpStr(s string, i int) (string, int) {
	j, maxj := i, len(s)
	for j < maxj && isOpChar(s[j]) {
		j++
	}
	return s[i:j], j
}
