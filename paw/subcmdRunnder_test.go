package main

import (
	//"fmt"
	"github.com/jlgerber/packwrap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubcmdRunner_New(t *testing.T) {
	nsc := NewSubcmdRunner()
	assert.NotNil(t, nsc)

	// test register
	nsc.RegisterSubcmd("foo", "foo it up", func(a *packwrap.ManifestLocator,
		b *packwrap.ManifestReaderFactory) error {
		println("foo")
		return nil
	})

	nsc.RegisterSubcmd("bar", "bar is the best", func(a *packwrap.ManifestLocator,
		b *packwrap.ManifestReaderFactory) error {
		println("bar")
		return nil
	})

	nsc.RegisterSubcmd("barbar", "barbar is the best", func(a *packwrap.ManifestLocator,
		b *packwrap.ManifestReaderFactory) error {
		println("barbar")
		return nil
	})

	// create manifest reader factory and register reader instances
	manifestReaderFactory := packwrap.NewManifestReaderFactory()
	manifestReaderFactory.AddReader("json", &packwrap.JsonManifestReader{})

	// create manifest locator with appropriate extensions
	manifestLocator := packwrap.NewManifestLocator(nsc.Keys())

	// register manifest locator and manifest reader factory with subcommand runner.
	nsc.RegisterManifestLocator(manifestLocator)
	nsc.RegisterManifestReaderFactory(manifestReaderFactory)

	assert.Equal(t, true, nsc.Has("foo"))
	assert.Equal(t, true, nsc.Has("bar"))
	assert.Equal(t, false, nsc.Has("bla"))

	// test get
	foo, ok := nsc.Get("foo")
	assert.Equal(t, ok, true)
	assert.NotNil(t, foo)
	//foo()

	bla, ok := nsc.Get("bla")
	assert.NotEqual(t, true, ok)
	assert.NotNil(t, bla)

	err := nsc.Run("foo")
	assert.Nil(t, err)

	err = nsc.Run("bla")
	assert.NotNil(t, err)

	desc := nsc.GetDesc("foo")
	assert.Equal(t, "foo it up", desc)

	maxlen := nsc.MaxNameLength()
	assert.Equal(t, 6, maxlen)

	keys := nsc.Keys()
	assert.Equal(t, 3, len(keys))
}
