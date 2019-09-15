package semver

import (
	"strings"
)

type ConstraintUnion uint8

const (
	ConstraintUnionOr ConstraintUnion = iota
	ConstraintUnionAnd
)

type Constraint struct {
	left  Checker
	right Checker
	un    ConstraintUnion
}

var _ Checker = (*Constraint)(nil)

func NewConstraint(s string) (*Constraint, error) {
	ors := strings.Split(s, "||")
	orConstr := make([]*Constraint, 0, len(ors))
	for _, or := range ors {
		ands := strings.Split(or, ",")
		andConstr := make([]*Constraint, 0, len(ands))
		for _, and := range ands {
			c, err := parseConstraint(and)
			if err != nil {
				return nil, err
			}
			andConstr = append(andConstr, c)
		}
		orConstr = append(orConstr, compact(andConstr, ConstraintUnionAnd))
	}
	return compact(orConstr, ConstraintUnionOr), nil
}

func (c *Constraint) Check(v *Version) bool {
	switch c.un {
	case ConstraintUnionAnd:
		return c.left.Check(v) && c.right.Check(v)
	case ConstraintUnionOr:
		return c.left.Check(v) || c.right.Check(v)
	}
	panic("should not happen")
}
