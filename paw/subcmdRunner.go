package main

import (
	"errors"
	"github.com/jlgerber/packwrap"
)

// SubcmdRunner - responsible for running subcommands, and providing subcommand descriptions.
// We register the subcommand functions, along with the ManifestLocator and ManifestReaderFactory here.
type SubcmdRunner struct {
	subcmd          map[string]func(*packwrap.ManifestLocator, *packwrap.ManifestReaderFactory) error
	desc            map[string]string
	manifestLocator *packwrap.ManifestLocator
	manifestRF      *packwrap.ManifestReaderFactory
}

// NewSubcmdRunner - constructor function
func NewSubcmdRunner() *SubcmdRunner {
	return &SubcmdRunner{subcmd: make(map[string]func(*packwrap.ManifestLocator,
		*packwrap.ManifestReaderFactory) error), desc: make(map[string]string)}
}

// RegisterSubcmd - register the subcommmand with the runner.
func (s *SubcmdRunner) RegisterSubcmd(name string, description string,
	subcmd func(*packwrap.ManifestLocator, *packwrap.ManifestReaderFactory) error) {
	s.subcmd[name] = subcmd
	s.desc[name] = description
}

// RegisterManifestLocator - set a pointer to the ManifestLocator instance, which is subsequently injected
// into the subcommand when run.
func (s *SubcmdRunner) RegisterManifestLocator(ml *packwrap.ManifestLocator) {
	s.manifestLocator = ml
}

// ReisterManifestReaderFactory - store a pointer to the ManifestReaderFactory instance, which is subsequently
// injected into the subcommand when it is run.
func (s *SubcmdRunner) RegisterManifestReaderFactory(mrf *packwrap.ManifestReaderFactory) {
	s.manifestRF = mrf
}

// Has - method to determine whether the runner has a subcommand with the supplied name.
func (s *SubcmdRunner) Has(name string) bool {
	_, ok := s.subcmd[name]
	return ok
}

// Get - return the subcommand with the matching name, as well as a boolean indicating success/failure
func (s *SubcmdRunner) Get(name string) (func(*packwrap.ManifestLocator, *packwrap.ManifestReaderFactory) error, bool) {
	val, ok := s.subcmd[name]
	if ok != true {
		return func(a *packwrap.ManifestLocator,
			b *packwrap.ManifestReaderFactory) error {
			return nil
		}, ok
	}
	return val, ok
}

// GetDesc - return the description string for the supplied name. If the name is not found, a generic message is
// returned indicating the failure.
func (s *SubcmdRunner) GetDesc(name string) string {
	d, ok := s.desc[name]
	if ok {
		return d
	}
	return "NO SUBCOMMAND NAMED " + name
}

// Run - execute the subcommand matching the supplied name. Return an error if unsuccessful, otherwise,
// return the returnvalue of the called subcommand.
func (s *SubcmdRunner) Run(name string) error {
	if s.manifestLocator == nil {
		return errors.New("ManifestLocator is nil. Was subcmdRunner.RegisterManifestLocator called?")
	}
	if s.manifestRF == nil {
		return errors.New("manifestRF is nil. was subcmdRunner.RegisterManifestReaderFactory called?")
	}
	value, ok := s.subcmd[name]
	if ok == false {
		return errors.New("subcommand runner - no subcommand named:" + name)
	}

	// call chosen function passing pointers to manifestLocator and manifestReaderFactory
	rv := value(s.manifestLocator, s.manifestRF)
	return rv
}

// MaxNameLength - return the length of the largest key in the runner
func (s *SubcmdRunner) MaxNameLength() int {
	maxlen := 0
	for key, _ := range s.desc {
		keylen := len(key)
		if keylen > maxlen {
			maxlen = keylen
		}
	}
	return maxlen
}

// Keys - return a slice of strings representing the name of each subcommand.
func (s *SubcmdRunner) Keys() []string {
	keys := make([]string, 0, len(s.desc))
	for k := range s.desc {
		keys = append(keys, k)
	}
	return keys
}
