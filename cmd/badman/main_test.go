package main_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/m-mizutani/badman"
	main "github.com/m-mizutani/badman/cmd/badman"
)

func TestDownloadAndDump(t *testing.T) {
	tmp, err := ioutil.TempFile("", "*.dat")
	require.NoError(t, err)
	tmp.Close()

	err = main.Handler([]string{"./badman", "dump", "-o", tmp.Name()})
	require.NoError(t, err)

	finfo, err := os.Stat(tmp.Name())
	require.NoError(t, err)
	assert.NotEqual(t, 0, finfo.Size())

	rfd, err := os.Open(tmp.Name())
	require.NoError(t, err)

	man := badman.New()
	assert.NoError(t, man.Load(rfd))

	os.Remove(tmp.Name())
}
