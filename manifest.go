package packwrap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

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
func GetManifestLocationFor(packageName, packageVersion string) (string, error) {
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

//--------------------------
// Types
//--------------------------
// Manifest - data structure
type Manifest struct {
	Name                  string
	Basepath              string // package install location
	Major                 uint16 // major version
	Minor                 uint16 // minor version
	Micro                 uint16 // micro version
	VersionTemplateString string // VersionStr to be optionally set by author
	Url                   string
	Environ               map[string]string
	// The following  fields are handled internally and should
	// not be set by the user.
	_version    *template.Template // rendered template. Handled internally
	_versionstr string             // private copy of the version string
}

//---------------------------
// methods
//---------------------------
func (m *Manifest) String() string {
	rs := fmt.Sprintf("%s-%s\n", m.Name, m.Version())
	for key, val := range m.Environ {
		rs += fmt.Sprintf("\t%s=%s\n", key, val)
	}
	return rs
}

// Version - return the version string, based on the Major Minor and
// Micro attributes, as well as the VersionTempalteString
func (m *Manifest) Version() string {
	if m._version == nil || m.VersionTemplateString != m._versionstr {
		if m.VersionTemplateString == "" {
			m.VersionTemplateString = VERSION_TEMPLATE_STRING
		}
		tmpl, err := template.New("version").Parse(m.VersionTemplateString)
		m._versionstr = m.VersionTemplateString
		if err != nil {
			panic(err)
		}
		m._version = tmpl
	}

	var outp bytes.Buffer
	if err := m._version.Execute(&outp, m); err != nil {
		panic(err)
	}

	return outp.String()
}

//----------------------------
// Constructors
//----------------------------
// NewManifestFromJsonByteSlice - given a byte slice encoded
// json document, return a new manifest by pointer.
func NewManifestFromJsonByteSlice(contents []byte) (*Manifest, error) {
	var manifest Manifest
	if err := json.Unmarshal(contents, &manifest); err != nil {
		return nil, err
	}
	// ah the joys of GC. No worries over cleaning up manifest
	return &manifest, nil
}

// NewManifestFromJsonFile return a pointer to a new manifest, and an error
// If the error is not nil, the *Manifest will be
func NewManifestFromJsonFile(jsonFile string) (*Manifest, error) {
	// open the file
	contents, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	return NewManifestFromJsonByteSlice(contents)
}
