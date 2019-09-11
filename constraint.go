package semver

type ConstraintOperator uint8

const (
	ConstraintAnd ConstraintOperator = iota
	ConstraintOr
)

type Constraint struct {
	left  Checker
	right Checker
	op    ConstraintOperator
}

var _ Checker = (*Constraint)(nil)

func NewConstraint(left Checker, op ConstraintOperator, right Checker) *Constraint {
	return &Constraint{
		left:  left,
		right: right,
		op:    op,
	}
}

func (c *Constraint) Check(v *Version) bool {
	switch c.op {
	case ConstraintAnd:
		return c.left.Check(v) && c.right.Check(v)
	case ConstraintOr:
		return c.left.Check(v) || c.right.Check(v)
	}
	panic("should not happen")
}
