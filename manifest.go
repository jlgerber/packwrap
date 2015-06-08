package packwrap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

// GetManifestSearchPathFor - returns the search path for
// the provided package
func GetManifestSearchPath() []string {

	manifestPath := DEFAULT_MANIFEST_LOCATION

	if manpath := os.Getenv(ENVVAR_MANIFESTPATH); manpath != "" {
		manifestPath = manpath + ":" + manifestPath
	}

	return strings.Split(manifestPath, ":")
}

// GetPackageVersions - given a package name, find all of the versions
// of the package and return them as a list
func GetPackageVersions(packageName string) []*PackageVersion {
	searchPath := GetManifestSearchPath()

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

//GetManifestFor - given the name of a package, return
// an error code and full path to the manifest assuming
// the returned error is nil.
func GetManifestLocationFor(packageName, packageVersion string) (string, error) {
	manifestPath := GetManifestSearchPath()

	for _, path := range manifestPath {
		manifest :=
			fmt.Sprintf("%s/%s/%s/manifest%s", path, packageName, packageVersion, Extension)
		//fmt.Println("searching", manifest)
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
	Schema                int
	Name                  string
	Basepath              string // package install location
	Major                 uint16 // major version
	Minor                 uint16 // minor version
	Micro                 uint16 // micro version
	VersionTemplateString string // VersionStr to be optionally set by author
	Url                   string
	Environ               []string
	// The following  fields are handled internally and should
	// not be set by the user.
	_version    *template.Template // rendered template. Handled internally
	_versionstr string             // private copy of the version string
}

// NewManifestFor - AlternateConstructor
func NewManifestFor(packageName, packageVersion string) (*Manifest, error) {
	manifestLocation, err := GetManifestLocationFor(packageName, packageVersion)
	if err != nil {
		return nil, err
	}

	manifest, err := NewManifestFromJsonFile(manifestLocation)

	if err != nil {
		return nil, err
	}
	return manifest, nil
}

//---------------------------
// methods
//---------------------------

// String - method returns a string represnetation of a Manifest
func (m *Manifest) String() string {
	rs := fmt.Sprintf("%s-%s\n", m.Name, m.Version())
	key := ""
	for idx, val := range m.Environ {
		if idx%2 == 0 {
			key = val
		} else {
			rs += fmt.Sprintf("\t%s=%s\n", key, val)
		}
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

// func (m *Manifest) ReplaceLocalVarsOld(s string) string {
// 	vals := strings.Fields(s)
// 	for i := range vals {
// 		if strings.HasPrefix(vals[i], "$$") {
// 			lookup := strings.Title(string(vals[i][2:]))
// 			newval, err := m.GetStringField(lookup)
// 			if err != nil {
// 				panic(err)
// 			}
// 			vals[i] = newval
// 		}
// 	}

// 	newstr := strings.Join(vals, " ")
// 	return newstr
// }

// replacement function
// ReplaceLocalVars - given a string with one or more
// local variables, prefixed by $$, replace them by
// looking up their value in the manifest and return
// the new string.
func (m *Manifest) ReplaceLocalVars(s string) string {

	ret, idxF, idxB, inVar := "", 0, 0, false

	for i, r := range s {
		if inVar == true {
			if string(r) == "/" || string(r) == ":" {

				lookup := strings.Title(strings.Trim(string(s[idxB+1:i]), "$"))
				newval, err := m.GetStringField(lookup)

				if err != nil {
					panic(err.Error() + lookup)
				}

				ret += newval

				idxB, idxF, inVar = i, i, false
			}
		} else if string(r) == "$" && string(s[i+1]) == "$" {
			ret = ret + s[idxB:idxF]
			idxB = idxF

			log.Debug("ReplaceLocalVars - idxB ", idxB, "idxF ", idxF, " s ", s)

			inVar = true
		} else {
			idxF = i
		}
	}
	// we have reached the end. account for the remaining data
	if idxB < len(s)-1 {
		log.Debug("ReplaceLocalVars - End of loop idxB", idxB, " last idx of s ", len(s)-1)
		// once last lookup
		if inVar == true {
			log.Debug("ReplaceLocalVars - End of loop inVar is true")

			lookup := strings.Title(strings.Trim(string(s[idxB+1:len(s)]), "$"))

			log.Debug("ReplaceLocalVars - Attempting to get field for >", lookup, "<")

			newval, err := m.GetStringField(lookup)
			if err != nil {
				panic(err.Error() + lookup)
			}
			ret = ret + newval
		} else {
			ret = ret + s[idxB:]
		}
	}

	return ret
}

// Setenv - set the environment of the current process based on the
// values in the Environ list of the Manifest.
func (m *Manifest) Setenv() error {
	var key, value string
	// had to switch from a dict to a slice in order to
	// preserve value
	for i, val := range m.Environ {
		if i%2 == 0 {
			key = val
			continue
		}
		value = val

		if strings.Contains(value, "$$") {
			value = m.ReplaceLocalVars(value)
		}
		// replace any shell variables defined so far with their value
		value = os.ExpandEnv(value)
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
		log.Debug("Setenv - Post", key, "=", value)
	}
	return nil
}

// Getenv - method retrieves the environment values from the Manifest,
// resolving shell and manifest variables....
func (m *Manifest) Getenv() []string {
	var key, value string
	var ret []string
	// had to switch from a dict to a slice in order to
	// preserve value
	for i, val := range m.Environ {
		if i%2 == 0 {
			key = val
			continue
		}
		value = val

		if strings.Contains(value, "$$") {
			value = m.ReplaceLocalVars(value)
		}
		// replace any shell variables defined so far with their value
		value = os.ExpandEnv(value)

		ret = append(ret, key+"="+value)
	}
	return ret
}

// GetStringField takes the name of a field and returns a string
// representing that fields value. This uses reflection, so asking
// for a non-string field will result in a panic....
// eg f = m.GetFieldValue("Foo")
func (m *Manifest) GetStringField(val string) (string, error) {
	r := reflect.ValueOf(m)
	v := reflect.Indirect(r)
	f := v.FieldByName(val)
	switch t := f.Kind(); {
	case t == reflect.String:
		return string(f.String()), nil

	case t == reflect.Uint16:
		return strconv.FormatUint(uint64(f.Uint()), 10), nil
	default:
		fmt.Println("kind of t", v.Kind())
	}
	return "", errors.New("GetStringField - Unable to convert field to type")
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
