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
}

type Constraint struct {
	left  Checker
	right Checker
	un    ConstraintUnion
}

var _ Checker = (*Constraint)(nil)

func NewConstraint(s string) (*Constraint, error) {
	//TODO
	return &Constraint{}, nil
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
