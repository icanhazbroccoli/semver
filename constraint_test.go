package semver

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	masterminds "github.com/Masterminds/semver"
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

func BenchmarkCheck(b *testing.B) {
	rand.Seed(time.Now().UTC().UnixNano())
	oc, _ := masterminds.NewConstraint("^10.20.30")
	nc := &Constraint{
		left:  NewGuard(newVersionUnsafe("10.20.30"), GuardGreaterOrEqual),
		right: NewGuard(newVersionUnsafe("11.0.0"), GuardLessThan),
		un:    ConstraintUnionAnd,
	}
	nVersions := 10000000
	type VerRes struct {
		Ver string
		Res bool
	}
	versions := make([]VerRes, 0, nVersions)
	for i := 0; i < nVersions; i++ {
		r := uint32(rand.Intn(0xFFFFFF))
		ver := fmt.Sprintf(
			"%d.%d.%d",
			(r>>16)&0xFF,
			(r>>8)&0xFF,
			r&0xFF,
		)
		sv, err := masterminds.NewVersion(ver)
		if err != nil {
			b.Fatal(err)
		}
		versions = append(versions, VerRes{Ver: ver, Res: oc.Check(sv)})
	}

	b.Run("icanhazbroccoli", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rv := versions[rand.Intn(len(versions))]
			v, err := NewVersion(rv.Ver)
			if err != nil {
				b.Fatal(err)
			}
			if c := nc.Check(v); c != rv.Res {
				b.Fatalf("icanhazbroccoli: mismatch for ver %s: got: %t, want: %t", rv.Ver, c, rv.Res)
			}
		}
	})

	b.Run("masterminds", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rv := versions[rand.Intn(len(versions))]
			v, err := masterminds.NewVersion(rv.Ver)
			if err != nil {
				b.Fatal(err)
			}
			if c := oc.Check(v); c != rv.Res {
				b.Fatalf("masterminds: mismatch for ver %s: got: %t, want: %t", rv.Ver, c, rv.Res)
			}
		}
	})
}
