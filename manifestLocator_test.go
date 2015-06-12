package packwrap

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestManifestLocator_NewManifestLocator(t *testing.T) {
	myLocator := NewManifestLocator([]string{"json", "toml", "yaml"})
	assert.NotNil(t, myLocator, "NewManifestLocator should return instance")
	ml := reflect.ValueOf(myLocator).Elem().Interface().(ManifestLocator)
	assert.NotNil(t, ml, "NewManifestLocator should return an instance of ManifestLocator")
}

func TestManifestLocator_GetManifestSearchPath(t *testing.T) {
	myLocator := NewManifestLocator([]string{"json", "toml", "yaml"})
	searchPath := myLocator.GetManifestSearchPath()
	assert.NotNil(t, searchPath, "Should not return Nil")
	assert.NotEmpty(t, searchPath, "GetManifestSearchPath should return non empty slice")

}

func TestManifestLocator_GetPackageVersions(t *testing.T) {
	myLocator := NewManifestLocator([]string{"json", "toml", "yaml"})
	versions := myLocator.GetPackageVersions("houdini")
	for _, version := range versions {
		assert.NotNil(t, reflect.ValueOf(version).Elem().Interface().(PackageVersion),
			"Should return list of PackageVersion")
		//fmt.Println(version.String())
	}
}

func TestManifestLocator_GetManifestLocationFor(t *testing.T) {
	myLocator := NewManifestLocator([]string{"json", "toml", "yaml"})
	location, err := myLocator.GetManifestLocationFor("maya", "2017.0.1")
	assert.Nil(t, err, "GetManifestLocationFor returned non-nil error")
	assert.NotNil(t, location)
	//fmt.Println(location)
}
