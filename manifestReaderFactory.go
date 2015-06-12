package packwrap

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type InvalidManifestReader struct{}

func (i InvalidManifestReader) NewManifestFromFile(f string) (*Manifest, error) {
	return nil, errors.New("InvalidManifestReader")
}

type ManifestReaderFactory struct {
	funcs map[string]ManifestReader //manifest reader is an interface.

}

func NewManifestReaderFactory() *ManifestReaderFactory {
	return &ManifestReaderFactory{funcs: make(map[string]ManifestReader)}

}

// AddReader - vaue is an interface, so you should treat it as a pointer.
// pass an address
func (m *ManifestReaderFactory) AddReader(name string, value ManifestReader) {
	m.funcs[name] = value
}

// HasReaderFor - tests wether the supplied format name is supported by the factory.
func (m *ManifestReaderFactory) HasReaderFor(name string) bool {
	_, ok := m.funcs[name]
	return ok
}

// GetReaderFor - get teh manifest reader for the provided name. If an invalid
// name is provided, return an instance of InvalidManifestReader.
func (m *ManifestReaderFactory) GetReaderFor(name string) ManifestReader {
	mapval, ok := m.funcs[name]
	if ok {
		return mapval
	}
	imr := InvalidManifestReader{}
	return imr
}

// NewManifestFor - Given the path to a manifest, return a pointer to an instance of
// Manifest and an error. If successful, the error will be nil, otherwise, the manifest
// pointer will be nil and the error will not be...
func (m *ManifestReaderFactory) NewManifestFor(path string) (*Manifest, error) {
	format := strings.Trim(filepath.Ext(path), ".")

	if m.HasReaderFor(format) == false {
		return nil, errors.New(fmt.Sprintf("ManifestReaderFactory.NewManifestFor - no reader for '%s'", format))
	}

	reader := m.GetReaderFor(format)

	return reader.NewManifestFromFile(path)
}
