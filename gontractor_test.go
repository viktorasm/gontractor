package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func requireFileExists(t *testing.T, file string) {
	_, err := os.Stat(file)
	require.NoError(t, err)
}

func TestE2E(t *testing.T) {
	g := NewGontractor()
	g.spec = "test-resources/swagger.yaml"
	g.apiOutFile = "test-output/testE2E/api/api.go_"
	g.serverOutFile = "test-output/testE2E/server_generated.go_"
	g.serverTemplate = "sample-templates/proprietary-api/server.tpl"
	g.Execute()

	requireFileExists(t, g.apiOutFile)
	requireFileExists(t, g.serverOutFile)
}

func TestGetPackageName(t *testing.T) {
	g := NewGontractor()

	noError := func(p string, err error) string {
		require.NoError(t, err)
		return p
	}

	require.Equal(t, "gontractor", noError(g.getPackageName("server.go")))
	require.Equal(t, "gontractor", noError(g.getPackageName("./server.go")))
	require.Equal(t, "foo", noError(g.getPackageName("foo/server.go")))
	require.Equal(t, "foo", noError(g.getPackageName("./foo/server.go")))
	require.Equal(t, "foo", noError(g.getPackageName("bar/foo/server.go")))
	require.Equal(t, "foo", noError(g.getPackageName("./bar/foo/server.go")))
}

func TestGetAbsolutePackagePath(t *testing.T) {
	g := NewGontractor()

	noError := func(p string, err error) string {
		require.NoError(t, err)
		return p
	}
	require.Equal(t, "github.com/viktorasm/gontractor", noError(g.getAbsolutePackagePath("server.go")))
	require.Equal(t, "github.com/viktorasm/gontractor", noError(g.getAbsolutePackagePath("./server.go")))
	require.Equal(t, "github.com/viktorasm/gontractor/bar/foo", noError(g.getAbsolutePackagePath("bar/foo/server.go")))
	require.Equal(t, "github.com/viktorasm/gontractor/bar/foo", noError(g.getAbsolutePackagePath("./bar/foo/server.go")))
}
