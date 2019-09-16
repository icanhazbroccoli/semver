# SemVer

SemVer is a library for a fast version match checking written in Golang. Can be used by package managers and other dependency-fetching clients dealing with semantic versioning where performance matters. The library works 5-7 times faster than similar known Golang implementations.

## Usage

```go
c, err := semver.NewConstraint("~>1.0.1")
if err != nil {
  ...
}
v, err := semver.NewVersion("1.0.2")
if err != nil {
  ...
}
if c.Check(v) {
  // version v belongs to the constraint-defined range
}
```

The library operates with 2 primitives: versions and constraints. A version defines a specific identifier, e.g.: `v1.0.1-beta.0`. A constraint defined an acceptable range of versions, e.g.: `~>1.0.1` means: `>=1.0.1 and < 1.1.0`. The library implements a fast checker for testing whether a given version belongs to the constraint-defined range.

## Implementation details and known limitations

The library stores version numbers in a single 32-bit unsigned integer value: 10 bit for every digit, therefore the amortised time cost of version comparison and increment/decrement operations is constant. In an optimistic scenario this happens in a single machine instruction. This introduces a limitation on the max version number: `1023.1023.1023`.
