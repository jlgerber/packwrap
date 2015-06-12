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
