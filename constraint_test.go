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

func TestParseConstraint(t *testing.T) {
	tests := []struct {
		Input        string
		ExpectConstr *Constraint
		ExpectErr    error
	}{
		{
			Input: "1.2.3",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, ""),
					GuardEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "1.2.*",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{1, 3, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "1.*",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{2, 0, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "*",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					&Version{base: 0x3FFFFFFF + 1},
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "*.*",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					&Version{base: 0x3FFFFFFF + 1},
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "*.*.*",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					&Version{base: 0x3FFFFFFF + 1},
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "=1.2.3-beta.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "=v1.2.3-beta.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "= v1.2.3-beta.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "!=1.2.3-beta.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardLessThan,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardGreaterThan,
				),
				un: ConstraintUnionOr,
			},
		},
		{
			Input: "!=1.2.*",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardLessThan,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{1, 3, 0}, ""),
					GuardGreaterThan,
				),
				un: ConstraintUnionOr,
			},
		},
		{
			Input: "!=1.2",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardLessThan,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{1, 3, 0}, ""),
					GuardGreaterThan,
				),
				un: ConstraintUnionOr,
			},
		},
		{
			Input: "!=1",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 0, 0}, ""),
					GuardLessThan,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{2, 0, 0}, ""),
					GuardGreaterThan,
				),
				un: ConstraintUnionOr,
			},
		},
		{
			Input: ">1.2.3-beta.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardGreaterThan,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: ">1.2.3",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, ""),
					GuardGreaterThan,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: ">1.2",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 3, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: ">1",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{2, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: ">=1.2.3-beta.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardGreaterOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: ">=1.2.3",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, ""),
					GuardGreaterOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: ">=1.2",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: ">=1",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "<1.2.3-beta.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardLessThan,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "<1.2.3",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, ""),
					GuardLessThan,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "<1.2",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardLessThan,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "<1",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 0, 0}, ""),
					GuardLessThan,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "<=1.2.3-beta.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardLessOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "<=1.2",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardLessOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "<=1",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 0, 0}, ""),
					GuardLessOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "~1.2.3-beta.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, "beta.0"),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{1, 3, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "~1.2.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{1, 3, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "~>1.2.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{1, 3, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "~1.2",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{1, 3, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "~1",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{2, 0, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "~>*",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "~>2.x.x",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{2, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{3, 0, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "^*",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: (*Guard)(nil),
				un:    ConstraintUnionOr,
			},
		},
		{
			Input: "^1.2.3",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 3}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{2, 0, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "^1.2",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 2, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{2, 0, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "^1",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{1, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{2, 0, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "^0.2.3",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 2, 3}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{0, 3, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "^0.2",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 2, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{0, 3, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "^0.0.3",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 0, 3}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{0, 0, 4}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "^0.0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{0, 1, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
		{
			Input: "^0",
			ExpectConstr: &Constraint{
				left: NewGuard(
					NewVersionRaw([]uint32{0, 0, 0}, ""),
					GuardGreaterOrEqual,
				),
				right: NewGuard(
					NewVersionRaw([]uint32{1, 0, 0}, ""),
					GuardLessThan,
				),
				un: ConstraintUnionAnd,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			c, err := parseConstraint(tt.Input)
			if !errorEqual(err, tt.ExpectErr) {
				t.Fatalf("unexpected error: got: %q, want: %q", err, tt.ExpectErr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(c, tt.ExpectConstr) {
				t.Fatalf("unexpected constraint: got: %#v, want: %#v", c, tt.ExpectConstr)
			}
		})
	}
}

type VerRes struct {
	Ver string
	Res bool
}

func BenchmarkStaticConstraintParseAndCheckCheck(b *testing.B) {
	rand.Seed(time.Now().UTC().UnixNano())
	oc, _ := masterminds.NewConstraint("^10.20.30")
	nc, _ := NewConstraint("^10.20.30")
	n := 10000
	versions := make([]VerRes, 0, n)
	for i := 0; i < n; i++ {
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

func BenchmarkParseConstraintOnCheck(b *testing.B) {
	rand.Seed(time.Now().UTC().UnixNano())
	n := 10000
	mc, _ := masterminds.NewConstraint("^10.20.30")
	versions := make([]VerRes, 0, n)
	for i := 0; i < n; i++ {
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
		versions = append(versions, VerRes{Ver: ver, Res: mc.Check(sv)})
	}

	b.Run("icanhazbroccoli", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			rv := versions[rand.Intn(len(versions))]
			nc, _ := NewConstraint("^10.20.30")
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
			oc, _ := masterminds.NewConstraint("^10.20.30")
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

func BenchmarkSimpleCompare(b *testing.B) {
	rand.Seed(time.Now().UTC().UnixNano())
	n := 100000
	nc, _ := NewConstraint("^10.20.30")
	oc, _ := masterminds.NewConstraint("^10.20.30")
	nVersions := make([]*Version, 0, n)
	oVersions := make([]*masterminds.Version, 0, n)
	for i := 0; i < n; i++ {
		r := uint32(rand.Intn(0xFFFFFF))
		ver := fmt.Sprintf(
			"%d.%d.%d",
			(r>>16)&0xFF,
			(r>>8)&0xFF,
			r&0xFF,
		)
		nv, _ := NewVersion(ver)
		nVersions = append(nVersions, nv)
		ov, _ := masterminds.NewVersion(ver)
		oVersions = append(oVersions, ov)
	}

	b.Run("icanhazbroccoli", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := nVersions[rand.Intn(len(nVersions))]
			nc.Check(v)
		}
	})

	b.Run("masterminds", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := oVersions[rand.Intn(len(oVersions))]
			oc.Check(v)
		}
	})
}
