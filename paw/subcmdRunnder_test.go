package main

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubcmdRunner_New(t *testing.T) {
	nsc := NewSubcmdRunner()
	assert.NotNil(t, nsc)

	// test register
	nsc.RegisterSubcmd("foo", "foo it up", func() error { println("foo"); return nil })
	nsc.RegisterSubcmd("bar", "bar is the best", func() error { println("bar"); return nil })
	nsc.RegisterSubcmd("barbar", "barbar is the best", func() error { println("barbar"); return nil })

	assert.Equal(t, true, nsc.Has("foo"))
	assert.Equal(t, true, nsc.Has("bar"))
	assert.Equal(t, false, nsc.Has("bla"))

	// test get
	foo, ok := nsc.Get("foo")
	assert.Equal(t, ok, true)
	assert.NotNil(t, foo)
	foo()

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
