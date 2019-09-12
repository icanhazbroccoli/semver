package semver

import (
	"reflect"
	"testing"
)

func TestConstraint(t *testing.T) {
	//c, err := NewConstraint(">=1.24.32-alpha.0, <2.0.0")
	v1, _ := NewVersion("1.24.32-alpha.0")
	v2, _ := NewVersion("2.0.0")
	c := &Constraint{
		left:  NewGuard(v1, GuardGreaterOrEqual),
		right: NewGuard(v2, GuardLessThan),
		un:    ConstraintUnionAnd,
	}
	//if err != nil {
	//	t.Fatal(err)
	//}

	tests := []struct {
		Name    string
		Version *Version
		Expect  bool
	}{
		{
			Name: "1.24.32-alpha.0",
			Version: &Version{
				base: (1 << 20) | (24 << 10) | 32,
				pre:  "alpha.0",
			},
			Expect: true,
		},
		{
			Name: "1.24.32-alpha.1",
			Version: &Version{
				base: (1 << 20) | (24 << 10) | 32,
				pre:  "alpha.1",
			},
			Expect: true,
		},
		{
			Name: "1.24.32",
			Version: &Version{
				base: (1 << 20) | (24 << 10) | 32,
			},
			Expect: true,
		},
		{
			Name: "1.24.31-alpha.0",
			Version: &Version{
				base: (1 << 20) | (24 << 10) | 31,
				pre:  "alpha.0",
			},
			Expect: false,
		},
		{
			Name: "1.23.32-alpha.0",
			Version: &Version{
				base: (1 << 20) | (23 << 10) | 32,
				pre:  "alpha.0",
			},
			Expect: false,
		},
		{
			Name: "0.24.32-alpha.0",
			Version: &Version{
				base: (0 << 20) | (24 << 10) | 32,
				pre:  "alpha.0",
			},
			Expect: false,
		},
		{
			Name: "2.0.0",
			Version: &Version{
				base: 2 << 20,
			},
			Expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if check := c.Check(tt.Version); check != tt.Expect {
				t.Fatalf("unexpected check result for %q: got: %t, want: %t", tt.Name, check, tt.Expect)
			}
		})
	}
}

func TestNewVersion(t *testing.T) {
	tests := []struct {
		Input     string
		ExpectErr error
		ExpectVer Version
	}{
		{
			Input: "1.2.3-beta.2",
			ExpectVer: Version{
				base: (1 << 20) | (2 << 10) | (3),
				pre:  "beta.2",
			},
		},
		{
			Input:     "0.0.0",
			ExpectVer: Version{},
		},
		{
			Input: "0.0.1",
			ExpectVer: Version{
				base: 1,
			},
		},
		{
			Input: "0.1.0",
			ExpectVer: Version{
				base: (1 << 10),
			},
		},
		{
			Input: "1.0.0",
			ExpectVer: Version{
				base: (1 << 20),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			ver, err := NewVersion(tt.Input)
			if !errorEqual(err, tt.ExpectErr) {
				t.Fatalf("unexpected error: got: %q, want: %q", err, tt.ExpectErr)
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(ver, &tt.ExpectVer) {
				t.Fatalf("unexpected version: got: %+v, want: %+v", *ver, tt.ExpectVer)
			}
		})
	}
}

func errorEqual(e1, e2 error) bool {
	if e1 == e2 {
		return true
	}
	if e1 != nil && e2 != nil {
		return e1.Error() == e2.Error()
	}
	return false
}
