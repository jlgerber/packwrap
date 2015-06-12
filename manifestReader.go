package packwrap

type ManifestReader interface {
	NewManifestFromFile(f string) (*Manifest, error)
}
