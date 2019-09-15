package semver

import "fmt"

type ConstraintOperator uint8

const (
	ConstraintOpTilde ConstraintOperator = iota
	ConstraintOpTildeOrEqual
	ConstraintOpNotEqual
	ConstraintOpGreaterThan
	ConstraintOpGreaterOrEqual
	ConstraintOpLessThan
	ConstraintOpLessOrEqual
	ConstraintOpCaret
)

var constraintOps map[string]ConstraintOperator

func init() {
	constraintOps = map[string]ConstraintOperator{
		"":   ConstraintOpTildeOrEqual,
		"=":  ConstraintOpTildeOrEqual,
		"!=": ConstraintOpNotEqual,
		">":  ConstraintOpGreaterThan,
		">=": ConstraintOpGreaterOrEqual,
		"=>": ConstraintOpGreaterOrEqual,
		"<":  ConstraintOpLessThan,
		"<=": ConstraintOpLessOrEqual,
		"=<": ConstraintOpLessOrEqual,
		"~":  ConstraintOpTilde,
		"~>": ConstraintOpTilde,
		"^":  ConstraintOpCaret,
	}
}

func parseConstraint(s string) (*Constraint, error) {
	var c *Constraint
	var left, right *Guard
	var un ConstraintUnion
	var err error

	i, maxi := 0, len(s)
	ds := make([]uint32, 0, 3)
	var wcds uint8
	var d, dix int
	i = skipTrailing(s, i)
	var op, pre string

	op, i = readOpStr(s, i)
	if _, ok := constraintOps[op]; !ok {
		return nil, fmt.Errorf("unrecognised constraint operator: %q", op)
	}
	i = skipTrailing(s, i)
	for i < maxi {
		if isNum(s[i]) {
			if dix >= 3 {
				goto Err
			}
			d, i = readNum(s, i)
			if i == -1 {
				goto Err
			}
			ds = append(ds, uint32(d))
			dix++
		} else if isStar(s[i]) {
			if dix >= 3 {
				goto Err
			}
			if w := uint8((1 << uint(2-dix))); w > wcds { // the most significant wildcard beats the rest
				wcds = w
			}
			ds = append(ds, 0)
			i++
			dix++
		} else {
			goto Err
		}
		if i < maxi && isDot(s[i]) {
			i++
			continue
		}
		if i < maxi && isDash(s[i]) {
			i++
			pre, _ = readStr(s, i)
			break
		}
	}

	left, right, un, err = genGuards(constraintOps[op], ds, wcds, pre)
	if err != nil {
		return nil, err
	}

	c = &Constraint{
		left:  left,
		right: right,
		un:    un,
	}
	return c, nil
Err:
	return nil, fmt.Errorf("failed to parse constraint %q around position %d", s, i)
}

func compact(cs []*Constraint, un ConstraintUnion) *Constraint {
	if len(cs) == 0 {
		return nil
	}
	ix := len(cs) - 1
	ptr := cs[ix]
	for ix > 0 {
		ix--
		ptr = &Constraint{
			left:  cs[ix],
			right: ptr,
			un:    un,
		}
	}
	return ptr
}

func expandRange(ds []uint32, wcds uint8, pre string) (*Version, *Version) {
	v := NewVersionRaw(ds, pre)
	switch wcds {
	case uint8(0):
		return v, v
	case uint8(1):
		return &Version{base: v.base & 0x3FFFFC00}, v.NextMinor()
	case uint8(2):
		return &Version{base: v.base & 0x3FF00000}, v.NextMajor()
	default:
		return &Version{base: 0}, &Version{base: 0x3FFFFFFF + 1}
	}
}

type guardGen func([]uint32, uint8, string) (*Guard, *Guard, ConstraintUnion, error)

func genGuards(op ConstraintOperator, ds []uint32, wcds uint8, pre string) (*Guard, *Guard, ConstraintUnion, error) {
	var gen guardGen
	switch op {
	case ConstraintOpTildeOrEqual:
		gen = genGuardTildeOrEqual
	case ConstraintOpNotEqual:
		gen = genGuardNotEqual
	}
	return gen(ds, wcds, pre)
}

func genGuardNotEqual(ds []uint32, wcds uint8, pre string) (*Guard, *Guard, ConstraintUnion, error) {
	v1, v2 := expandRange(ds, wcds, pre)
	return NewGuard(v1, GuardLessThan), NewGuard(v2, GuardGreaterThan), ConstraintUnionOr, nil
}

func genGuardTildeOrEqual(ds []uint32, wcds uint8, pre string) (*Guard, *Guard, ConstraintUnion, error) {
	v1, v2 := expandRange(ds, wcds, pre)
	if v1 == v2 {
		return NewGuard(v1, GuardEqual), nil, ConstraintUnionOr, nil
	}
	return NewGuard(v1, GuardGreaterOrEqual), NewGuard(v2, GuardLessThan), ConstraintUnionAnd, nil
}

func genGuardGreaterThan(ds []uint32, wcds uint8, pre string) (*Guard, *Guard, ConstraintUnion, error) {
	_, v2 := expandRange(ds, wcds, pre)
	return NewGuard(v2, GuardGreaterThan), nil, ConstraintUnionOr, nil
}

func genGuardGreaterOrEqual(ds []uint32, wcds uint8, pre string) (*Guard, *Guard, ConstraintUnion, error) {
	_, v2 := expandRange(ds, wcds, pre)
	return NewGuard(v2, GuardGreaterOrEqual), nil, ConstraintUnionOr, nil
}

func genGuardLessThan(ds []uint32, wcds uint8, pre string) (*Guard, *Guard, ConstraintUnion, error) {
	v1, _ := expandRange(ds, wcds, pre)
	return NewGuard(v1, GuardLessThan), nil, ConstraintUnionOr, nil
}

func genGuardLessOrEqual(ds []uint32, wcds uint8, pre string) (*Guard, *Guard, ConstraintUnion, error) {
	v1, _ := expandRange(ds, wcds, pre)
	return NewGuard(v1, GuardLessOrEqual), nil, ConstraintUnionOr, nil
}

func genGuardTilde(ds []uint32, wcds uint8, pre string) (*Guard, *Guard, ConstraintUnion, error) {
	v1, v2 := expandRange(ds, wcds, pre)
	return NewGuard(v1, GuardGreaterOrEqual), NewGuard(v2, GuardLessThan), ConstraintUnionAnd, nil
}

//func genGuardCaret(ds [3]uint32), wcds uint8, pre string) (*Guard, *Guard, ConstraintUnion, error) {
//
//}
