package comparer

import (
	"strings"

	"golang.org/x/xerrors"

	"github.com/khulnasoft-lab/goversion/pkg/version"
	dbTypes "github.com/khulnasoft-lab/vul-db/pkg/types"
	"github.com/khulnasoft-lab/vul/pkg/log"
)

// Comparer is an interface for version comparison
type Comparer interface {
	IsVulnerable(currentVersion string, advisory dbTypes.Advisory) bool
}

type matchVersion func(currentVersion, constraint string) (bool, error)

// IsVulnerable checks if the package version is vulnerable to the advisory.
func IsVulnerable(pkgVer string, advisory dbTypes.Advisory, match matchVersion) bool {
	// If one of vulnerable/patched versions is empty, we should detect it anyway.
	for _, v := range append(advisory.VulnerableVersions, advisory.PatchedVersions...) {
		if v == "" {
			return true
		}
	}
	var matched bool
	var err error

	if len(advisory.VulnerableVersions) != 0 {
		matched, err = match(pkgVer, strings.Join(advisory.VulnerableVersions, " || "))
		if err != nil {
			log.Logger.Warn(err)
			return false
		} else if !matched {
			// the version is not vulnerable
			return false
		}
	}

	secureVersions := append(advisory.PatchedVersions, advisory.UnaffectedVersions...)
	if len(secureVersions) == 0 {
		// the version matches vulnerable versions and patched/unaffected versions are not provided
		// or all values are empty
		return matched
	}

	matched, err = match(pkgVer, strings.Join(secureVersions, " || "))
	if err != nil {
		log.Logger.Warn(err)
		return false
	}
	return !matched
}

// GenericComparer represents a comparer for semver-like versioning
type GenericComparer struct{}

// IsVulnerable checks if the package version is vulnerable to the advisory.
func (v GenericComparer) IsVulnerable(ver string, advisory dbTypes.Advisory) bool {
	return IsVulnerable(ver, advisory, v.matchVersion)
}

// matchVersion checks if the package version satisfies the given constraint.
func (v GenericComparer) matchVersion(currentVersion, constraint string) (bool, error) {
	ver, err := version.Parse(currentVersion)
	if err != nil {
		return false, xerrors.Errorf("version error (%s): %s", currentVersion, err)
	}

	c, err := version.NewConstraints(constraint)
	if err != nil {
		return false, xerrors.Errorf("constraint error (%s): %s", currentVersion, err)
	}

	return c.Check(ver), nil
}
