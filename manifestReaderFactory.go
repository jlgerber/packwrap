package packwrap

import "errors"

type InvalidManifestReader struct{}

func (i InvalidManifestReader) NewManifestFromFile(f string) (*Manifest, error) {
	return nil, errors.New("InvalidManifestReader")
}

type ManifestReaderFactory struct {
	funcs map[string]ManifestReader
}

func NewManifestReaderFactory() *ManifestReaderFactory {
	return &ManifestReaderFactory{funcs: make(map[string]ManifestReader)}

}

func (m *ManifestReaderFactory) AddReader(name string, value ManifestReader) {
	m.funcs[name] = value
}

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
