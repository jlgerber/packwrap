package packwrap

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestManifestReaderFactory_AddReader(t *testing.T) {
	//create manifest
	jmr := JsonManifestReader{}
	myFactory := NewManifestReaderFactory()

	myFactory.AddReader("json", &jmr)
	assert.Equal(t, true, myFactory.HasReaderFor("json"),
		"myFactory.HasReaderFor(json) failed")

}

func TestManifestReaderFactory_GetReader(t *testing.T) {
	//create manifest
	jmr := JsonManifestReader{}
	myFactory := NewManifestReaderFactory()

	myFactory.AddReader("json", &jmr)
	assert.Equal(t, true, myFactory.HasReaderFor("json"), "myFactory.HasReaderFor(json) failed")

	myReader := myFactory.GetReaderFor("json")
	jt := reflect.ValueOf(myReader).Elem().Interface().(JsonManifestReader)
	assert.NotNil(t, jt)
	//println("myReader", myReader)
}

func TestManifestReaderFactory_GetManifest(t *testing.T) {
	//create manifest
	jmr := JsonManifestReader{}
	myFactory := NewManifestReaderFactory()

	myFactory.AddReader("json", &jmr)
	assert.Equal(t, true, myFactory.HasReaderFor("json"), "myFactory.HasReaderFor(json) failed")

	myManifest, err := myFactory.NewManifestFor("/packages/manifest/houdini/14.0.335/manifest.json")
	assert.Nil(t, err, "unable to read manifest")
	jt := reflect.ValueOf(myManifest).Elem().Interface().(Manifest)
	assert.NotNil(t, jt)
	//println("myReader", myReader)
}
