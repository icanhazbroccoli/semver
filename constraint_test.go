package semver

import (
	"reflect"
	"testing"
)

// 1-liner version initialization, returns nil on error
func newVersionUnsafe(s string) *Version {
	if v, err := NewVersion(s); err != nil {
		return nil
	} else {
		return v
	}
}

func TestCompact(t *testing.T) {
	tests := []struct {
		Name   string
		Input  []*Constraint
		Union  ConstraintUnion
		Expect *Constraint
	}{
		{
			Name: "Single constraint",
			Input: []*Constraint{
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("1.2.3"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("2.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
			},
			Union: ConstraintUnionAnd,
			Expect: &Constraint{
				left: NewGuard(
					newVersionUnsafe("1.2.3"),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					newVersionUnsafe("2.0.0"),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},

		{
			Name: "2 constraints",
			Input: []*Constraint{
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("1.2.3"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("2.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("4.5.6"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("5.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
			},
			Union: ConstraintUnionOr,
			Expect: &Constraint{
				left: &Constraint{
					left: NewGuard(
						newVersionUnsafe("1.2.3"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("2.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				right: &Constraint{
					left: NewGuard(
						newVersionUnsafe("4.5.6"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("5.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				un: ConstraintUnionOr,
			},
		},

		{
			Name: "3 constraints",
			Input: []*Constraint{
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("1.2.3"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("2.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("4.5.6"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("5.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("7.8.9"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("8.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
			},
			Union: ConstraintUnionOr,
			Expect: &Constraint{
				left: &Constraint{
					left: NewGuard(
						newVersionUnsafe("1.2.3"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("2.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				right: &Constraint{
					left: &Constraint{
						left: NewGuard(
							newVersionUnsafe("4.5.6"),
							GuardGreaterOrEqual,
						),
						right: NewGuard(
							newVersionUnsafe("5.0.0"),
							GuardLessThan,
						),
						un: ConstraintUnionAnd,
					},
					right: &Constraint{
						left: NewGuard(
							newVersionUnsafe("7.8.9"),
							GuardGreaterOrEqual,
						),
						right: NewGuard(
							newVersionUnsafe("8.0.0"),
							GuardLessThan,
						),
						un: ConstraintUnionAnd,
					},
					un: ConstraintUnionOr,
				},
				un: ConstraintUnionOr,
			},
		},

		{
			Name: "4 constraints",
			Input: []*Constraint{
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("1.2.3"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("2.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("4.5.6"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("5.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("7.8.9"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("8.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				&Constraint{
					left: NewGuard(
						newVersionUnsafe("10.11.12"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("11.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
			},
			Union: ConstraintUnionOr,
			Expect: &Constraint{
				left: &Constraint{
					left: NewGuard(
						newVersionUnsafe("1.2.3"),
						GuardGreaterOrEqual,
					),
					right: NewGuard(
						newVersionUnsafe("2.0.0"),
						GuardLessThan,
					),
					un: ConstraintUnionAnd,
				},
				right: &Constraint{
					left: &Constraint{
						left: NewGuard(
							newVersionUnsafe("4.5.6"),
							GuardGreaterOrEqual,
						),
						right: NewGuard(
							newVersionUnsafe("5.0.0"),
							GuardLessThan,
						),
						un: ConstraintUnionAnd,
					},
					right: &Constraint{
						left: &Constraint{
							left: NewGuard(
								newVersionUnsafe("7.8.9"),
								GuardGreaterOrEqual,
							),
							right: NewGuard(
								newVersionUnsafe("8.0.0"),
								GuardLessThan,
							),
							un: ConstraintUnionAnd,
						},
						right: &Constraint{
							left: NewGuard(
								newVersionUnsafe("10.11.12"),
								GuardGreaterOrEqual,
							),
							right: NewGuard(
								newVersionUnsafe("11.0.0"),
								GuardLessThan,
							),
							un: ConstraintUnionAnd,
						},
						un: ConstraintUnionOr,
					},
					un: ConstraintUnionOr,
				},
				un: ConstraintUnionOr,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			cc := compact(tt.Input, tt.Union)
			if !reflect.DeepEqual(cc, tt.Expect) {
				t.Fatalf("unexpected constraint: got: %+v, want: %+v", *cc, *tt.Expect)
			}
		})
	}
}
