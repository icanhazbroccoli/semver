# SemVer

SemVer is a library for a fast version match checking written in Golang. Can be used by package managers and other dependency-fetching clients dealing with semantic versioning where performance matters. The library works several times faster than similar known Golang implementations (see Benchmarks section).

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

## Benchmarks

In the benchmarks the library performance is compared against [Masterminds/semver](https://github.com/Masterminds/semver). This library is a very comprehensive tool to operate with SemVer constraints and versions.

The CPU is:
```
Model: Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz
# of cores: 8
Cache size: 8192 KB
Enabled flags: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc art arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc cpuid aperfmperf tsc_known_freq pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch cpuid_fault epb invpcid_single pti ssbd ibrs ibpb stibp tpr_shadow vnmi flexpriority ept vpid fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx rdseed adx smap clflushopt intel_pt xsaveopt xsavec xgetbv1 xsaves dtherm ida arat pln pts hwp hwp_notify hwp_act_window hwp_epp md_clear flush_l1d
```

Results:
```
goos: linux
goarch: amd64
pkg: sandbox/semver
BenchmarkStaticConstraintParseAndCheckCheck/icanhazbroccoli-8            4695884               348 ns/op
BenchmarkStaticConstraintParseAndCheckCheck/masterminds-8                 405423              3313 ns/op
BenchmarkParseConstraintOnCheck/icanhazbroccoli-8                         770362              2106 ns/op
BenchmarkParseConstraintOnCheck/masterminds-8                              84526             13316 ns/op
BenchmarkSimpleCompare/icanhazbroccoli-8                                 9843217               112 ns/op
BenchmarkSimpleCompare/masterminds-8                                     5460199               187 ns/op
```

In the first pair of tests a caret constraint was compiled once and compared
against 10k versions parsed on every iteration.

In the second test we parsed caret constraint and a version on every check 10k
times.

In the third benchmark we pre-compiled versions and the caret constraint and
benchmarked pure Check() call.

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
