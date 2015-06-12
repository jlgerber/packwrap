package main

import (
	"errors"
)

type SubcmdRunner struct {
	subcmd map[string]func() error
	desc   map[string]string
}

func NewSubcmdRunner() *SubcmdRunner {
	return &SubcmdRunner{subcmd: make(map[string]func() error), desc: make(map[string]string)}
}

// RegisterSubcmd - register the subcommmand with the runner.
func (s *SubcmdRunner) RegisterSubcmd(name string, description string, subcmd func() error) {
	s.subcmd[name] = subcmd
	s.desc[name] = description
}

func (s *SubcmdRunner) Has(name string) bool {
	_, ok := s.subcmd[name]
	return ok
}

func (s *SubcmdRunner) Get(name string) (func() error, bool) {
	val, ok := s.subcmd[name]
	if ok != true {
		return func() error { return nil }, ok
	}
	return val, ok
}

func (s *SubcmdRunner) GetDesc(name string) string {
	d, ok := s.desc[name]
	if ok {
		return d
	}
	return "NO SUBCOMMAND NAMED " + name
}

func (s *SubcmdRunner) Run(name string) error {
	value, ok := s.subcmd[name]
	if ok == false {
		return errors.New("subcommand runner - no subcommand named:" + name)
	}

	rv := value()
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

func (s *SubcmdRunner) Keys() []string {
	keys := make([]string, 0, len(s.desc))
	for k := range s.desc {
		keys = append(keys, k)
	}
	return keys
}
