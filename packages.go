package packwrap

import (
	"fmt"
	"io/ioutil"
)

// GetPackageList - function to return a list of packages which exist in one
// of the manifest directories.
func GetPackageList() []string {
	rl := []string{}
	searchPath := GetManifestSearchPath()
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
