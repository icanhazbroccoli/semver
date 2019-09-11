package semver

type GuardEquality uint8

const (
	GuardEqual GuardEquality = iota
	GuardGreaterThan
	GuardGreaterOrEqual
	GuardLessThan
	GuardLessOrEqual
)

type Guard struct {
	ver *Version
	op  GuardEquality
}

func NewGuard(ver *Version, op GuardEquality) *Guard {
	return &Guard{
		ver: ver,
		op:  op,
	}
}

func (g *Guard) Check(v *Version) bool {
	eq := g.ver.Equal(v)
	less := !eq && v.Less(g.ver)
	switch g.op {
	case GuardEqual:
		return eq
	case GuardGreaterThan:
		return !less
	case GuardGreaterOrEqual:
		return eq || !less
	case GuardLessThan:
		return less
	case GuardLessOrEqual:
		return eq || less
	}
	panic("should not happen either")
}

var _ Checker = (*Guard)(nil)
