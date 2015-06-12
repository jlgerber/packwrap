package packwrap

import (
	"encoding/json"
	"io/ioutil"
)

type JsonManifestReader struct{}

// NewManifestFromJsonByteSlice - given a byte slice encoded
// json document, return a new manifest by pointer.
func (jmr JsonManifestReader) NewManifestFromByteSlice(contents []byte) (*Manifest, error) {
	var manifest Manifest
	if err := json.Unmarshal(contents, &manifest); err != nil {
		return nil, err
	}
	// ah the joys of GC. No worries over cleaning up manifest
	return &manifest, nil
}
func (jmr JsonManifestReader) NewManifestFromFile(f string) (*Manifest, error) {

	contents, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	return jmr.NewManifestFromByteSlice(contents)
}
