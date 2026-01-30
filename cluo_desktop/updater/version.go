package updater

import (
	"fmt"
	"strconv"
	"strings"
)

// SemVer represents a semantic version.
type SemVer struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
}

// ParseSemVer parses a semantic version string (e.g., "1.2.3", "1.2.3-beta").
func ParseSemVer(version string) (SemVer, error) {
	v := SemVer{}

	// Remove leading 'v' if present
	version = strings.TrimPrefix(version, "v")

	// Split prerelease suffix
	parts := strings.SplitN(version, "-", 2)
	if len(parts) == 2 {
		v.Prerelease = parts[1]
	}

	// Parse major.minor.patch
	versionParts := strings.Split(parts[0], ".")
	if len(versionParts) < 1 || len(versionParts) > 3 {
		return v, fmt.Errorf("invalid version format: %s", version)
	}

	var err error
	v.Major, err = strconv.Atoi(versionParts[0])
	if err != nil {
		return v, fmt.Errorf("invalid major version: %s", versionParts[0])
	}

	if len(versionParts) >= 2 {
		v.Minor, err = strconv.Atoi(versionParts[1])
		if err != nil {
			return v, fmt.Errorf("invalid minor version: %s", versionParts[1])
		}
	}

	if len(versionParts) >= 3 {
		v.Patch, err = strconv.Atoi(versionParts[2])
		if err != nil {
			return v, fmt.Errorf("invalid patch version: %s", versionParts[2])
		}
	}

	return v, nil
}

// String returns the string representation of the version.
func (v SemVer) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		s += "-" + v.Prerelease
	}
	return s
}

// Compare compares two versions.
// Returns -1 if v < other, 0 if v == other, 1 if v > other.
func (v SemVer) Compare(other SemVer) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}

	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}

	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}

	// Prerelease comparison: no prerelease > any prerelease
	// e.g., 1.0.0 > 1.0.0-beta
	if v.Prerelease == "" && other.Prerelease != "" {
		return 1
	}
	if v.Prerelease != "" && other.Prerelease == "" {
		return -1
	}
	if v.Prerelease < other.Prerelease {
		return -1
	}
	if v.Prerelease > other.Prerelease {
		return 1
	}

	return 0
}

// IsNewerThan returns true if v is newer than other.
func (v SemVer) IsNewerThan(other SemVer) bool {
	return v.Compare(other) > 0
}
