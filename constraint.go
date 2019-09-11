package semver

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

type Constraint struct {
	left  Checker
	right Checker
	un    ConstraintUnion
}

var _ Checker = (*Constraint)(nil)

func NewConstraint(left Checker, un ConstraintUnion, right Checker) *Constraint {
	return &Constraint{
		left:  left,
		right: right,
		un:    un,
	}
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
