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

## Understanding SemVer constraints

Every SemVer constraint can be presented with either an open range
(single-ended) or a closed range (double-ended). An example of an open range is
a constraint `>0.1.2`, or: `<=1.2.*`. A closed range can look like: `>=1.2.3,
<2.0.0`.

Semver constraints introduce a range of operators: extra helpers that make these
ranges look clear, descriptive and self-explanatory. In fact, SemVer can always
always be expanded to a range. Having operators can be seen as syntax sugar.

There are several SemVer operators:

|Operator|Example  |Equivalent Range   |Explanation|
|--------|---------|-------------------|-----------|
|   `=`  |`=1.2.3` |`=1.2.3`           |A single specific version. The equal sign is optional, as well as the leading v|
|        |`=v1.2`  |`>=1.2.0, <1.3.0`  |1.2 defines a family of 1.2s, where it ranges between 1.2.0(incl) and 1.3.0(excl)
|        |`1.2.x`  |`>=1.2.0, <1.3.0`  |An equivalent to the example above but using a wildcard. `*`, `x` and `X` are total equivalents|
|  `!=`  |`!=1.2.3`|`!=1.2.3`          |Anything but this specific version|
|        |`!=1.2`  |`<1.2.0 || >=1.3.0`|1.2 is a range family, negating the range gives 2 open ranges|
|  `>`   |`>1.2.3` |`>1.2.3`           |Trivial non-inclusive range|
|        |`>1.2`   |`>=1.3.0`          |1.2s is between 1.2.0(incl) and 1.3.0(excl). It's 1.3.0(incl) to the right from this range|
|        |`>1`     |`>=2.0.0`          |1s is even a bigger range: it expands to the next major version|
|`>=`,`=<`|`>=1.2.3`|`>=1.2.3`          |Trivial inclusive range|
|        |`>=1.2`  |`>=1.2.0`          |This time range between 1.2.0 and 1.3.0 is included and expands in a single-ended range|
|        |`>=0`    |`>=0.0.0`          |Any version satisfies this constraint|
|  `<`   |`<1.2.3` |`<1.2.3`           |Trivial non-inclusive range|
|        |`<1.2`   |`<1.2.0`           |A ranged version expands in a trivial non-inclusive range|
|        |`<0`     |`<0.0.0`           |No version satisifies this constraint, oppsoite to `>=0.0.0`|
|`<=`,`=<`|`<=1.2.3`|`<=1.2.3`          |Trivial inclusive range|
|`~`,`~>`|`~1.2.3` |`>=1.2.3, <1.3.0`  |When either a minor version and a patch or just a minor version specified, tilde expands to the next minor release|
|        |`~1.2`   |`>=1.2.0, <1.3.0`  |           |
|        |`~1`     |`>=1.0.0, <2.0.0`  |If only a major specified, tilde expands to the next major|
|        |`~*`     |`>=0.0.0`          |All-matching constraint|
|  `^`   |`^1.2.3` |`>=1.2.3, <2.0.0`  |Caret has an extra contextual dependency: it changes it's behavior depending on whether major and minor versions are zero or not. If major is non-zero, caret expands to the next major|
|        |`^1.2`   |`>=1.2.0, <2.0.0`  |           |
|        |`^1`     |`>=1.0.0, <2.0.0`  |           |
|        |`^*`     |`>=0.0.0`          |This one is special and expands to an all-match constraint|
|        |`^0.2.3` |`>=0.2.3, <0.3.0`  |Major is 0 but minor is not, therefore caret expands the version to the next minor|
|        |`^0.0.3` |`>=0.0.3, <0.0.4`  |Major and minor are 0, the constraint is expanded to the next patch|
|        |`^0.0`   |`>=0.0.0, <0.1.0`  |Major and minor are 0, but the constraint expands to a range, therefore it bounded by the next minor|
|        |`^0`     |`>=0.0.0, <1.0.0`  |Major is 0, but the constraint expands to a range, therefore it bounded by the next major|


## Implementation details and known limitations

The library stores version numbers in a single 32-bit unsigned integer value: 10 bit for every digit, therefore the amortised time cost of version comparison and increment/decrement operations is constant. In an optimistic scenario this happens in a single machine instruction. This introduces a limitation on the max version number: `1023.1023.1023`.
