package semver

import (
	"fmt"
	"regexp"
	"strings"
)

type ConstraintUnion uint8

const (
	ConstraintUnionAnd ConstraintUnion = iota
	ConstraintUnionOr
)

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

const cvRegex string = `v?([0-9|x|X|\*]+)(\.[0-9|x|X|\*]+)?(\.[0-9|x|X|\*]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

var constraintOps map[string]ConstraintOperator

var constraintRegex *regexp.Regexp
var findConstraintRegex *regexp.Regexp
var validConstraintRegex *regexp.Regexp

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

	ops := make([]string, 0, len(constraintOps))
	for k := range constraintOps {
		ops = append(ops, regexp.QuoteMeta(k))
	}

	constraintRegex = regexp.MustCompile(fmt.Sprintf(
		`^\s*(%s)\s*(%s)\s*$`,
		strings.Join(ops, "|"),
		cvRegex))

	findConstraintRegex = regexp.MustCompile(fmt.Sprintf(
		`(%s)\s*(%s)`,
		strings.Join(ops, "|"),
		cvRegex))

	validConstraintRegex = regexp.MustCompile(fmt.Sprintf(
		`^(\s*(%s)\s*(%s)\s*\,?)+$`,
		strings.Join(ops, "|"),
		cvRegex))
}

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

func parseConstraint(s string) (*Constraint, error) {
	return nil, fmt.Errorf("not implemented")
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
