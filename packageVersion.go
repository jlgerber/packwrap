package packwrap

import (
	"fmt"
)

// PackageVersion - struct which holds information about a package
type PackageVersion struct {
	Name     string // The name of the package
	Version  string // the version string
	Location string // the location of the version
}

// NewPackageVersion - constructor which returns a pointer to a new
// PackageVersion
func NewPackageVersion(name, version, location string) *PackageVersion {
	pv := PackageVersion{name, version, location}
	return &pv
}

// String - method which returns a string representation of a PackageVersion
func (p *PackageVersion) String() string {
	return fmt.Sprintf("%s-%s %s", p.Name, p.Version, p.Location)
}
