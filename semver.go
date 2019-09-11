package semver

type Checker interface {
	Check(*Version) bool
}
