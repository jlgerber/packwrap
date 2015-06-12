package packwrap

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type ManifestLocator struct {
	extensions []string
}

// NewManifestLocator - constructor function. Given a string slice containing valid extensions, initialize
// the ManifestLocator.
func NewManifestLocator(extensions []string) *ManifestLocator {
	return &ManifestLocator{extensions: extensions}
}

func (m *ManifestLocator) GetPackageList() []string {
	rl := []string{}
	searchPath := m.GetManifestSearchPath()
	for _, path := range searchPath {
		info, err := ioutil.ReadDir(path)
		if err != nil {
			log.Debug(err)
			continue
		}
		for _, pack := range info {
			if string(pack.Name()[0]) == "." {
				continue
			}
			rl = append(rl, fmt.Sprintf("%s    %s", pack.Name(), path))
		}
	}
	return rl
}

// GetManifestSearchPathFor - returns the search path for
// the provided package
func (m *ManifestLocator) GetManifestSearchPath() []string {

	manifestPath := DEFAULT_MANIFEST_LOCATION

	if manpath := os.Getenv(ENVVAR_MANIFESTPATH); manpath != "" {
		manifestPath = manpath + ":" + manifestPath
	}

	return strings.Split(manifestPath, ":")
}

// GetPackageVersions - given a package name, find all of the versions
// of the package and return them as a list
func (m *ManifestLocator) GetPackageVersions(packageName string) []*PackageVersion {
	searchPath := m.GetManifestSearchPath()

	versions := make([]*PackageVersion, 0)

	for _, path := range searchPath {
		packagePath := path + "/" + packageName

		info, err := ioutil.ReadDir(packagePath)
		if err != nil {
			log.Debug(err)
			continue
		}
		for _, version := range info {
			if string(version.Name()[0]) == "." {
				continue
			}
			versions = append(versions,
				NewPackageVersion(packageName, version.Name(), path))
		}
	}
	return versions
}

// GetManifestFor - given the name of a package, return
// an error code and full path to the manifest assuming
// the returned error is nil.
func (m *ManifestLocator) GetManifestLocationFor(packageName, packageVersion string) (string, error) {
	manifestPath := m.GetManifestSearchPath()

	for _, path := range manifestPath {
		for _, extension := range m.extensions {
			manifest := fmt.Sprintf("%s/%s/%s/manifest.%s", path, packageName, packageVersion, extension)
			//fmt.Println("searching", manifest)
			log.Debugf("GetManifestLocationFor(%s, %s) - searching %s", packageName, packageVersion, manifest)
			if _, err := os.Stat(manifest); err == nil {
				return manifest, nil
			}
		}
	}

	return "", errors.New("Unable to find manifest.")
}
