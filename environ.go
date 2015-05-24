package packwrap

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const ManifestPath string = "/Users/jonathangerber/manifest"

// GetManifestSearchPathFor - returns the search path for
// the provided package
func GetManifestSearchPathFor(packageName, packageVersion string) string {
	manifestPath := ManifestPath

	if manpath := os.Getenv(Envvar_manifestPath); manpath != "" {
		manifestPath = manpath + ":" + manifestPath
	}

	return manifestPath
}

//GetManifestFor - given the name of a package, return
// an error code and full path to the manifest assuming
// the returned error is nil.
func GetManifestFor(packageName string, packageVersion string) (string, error) {
	manifestPath := GetManifestSearchPathFor(packageName, packageVersion)

	for _, path := range strings.Split(manifestPath, ":") {
		manifest := path +
			fmt.Sprintf("/%s/%s/manifest%s", packageName, packageVersion, Extension)

		if _, err := os.Stat(manifest); err == nil {
			return manifest, nil
		}
	}

	return "", errors.New("Unable to find manifest.")
}
