package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
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

	validate := func(expectedPackageName string, packageFile string) {
		name, err := g.getPackageName(packageFile)
		require.NoError(t, err)
		require.Equal(t, name, expectedPackageName)
	}

	validate("gontractor", "server.go")
	validate("gontractor", "./server.go")
	validate("foo", "foo/server.go")
	validate("foo", "./foo/server.go")
	validate("foo", "bar/foo/server.go")
	validate("foo", "./bar/foo/server.go")
}

func TestGetAbsolutePackagePath(t *testing.T) {
	g := NewGontractor()

	validate := func(expectedPackagePath string, relativeFileName string) {
		relativeFileName = filepath.FromSlash(relativeFileName)
		path, err := g.getAbsolutePackagePath(relativeFileName)
		require.NoError(t, err)
		require.Equal(t, expectedPackagePath, path)
	}

	// handle relative paths: for tests this will equate to this same package
	validate("github.com/viktorasm/gontractor", "server.go")
	validate("github.com/viktorasm/gontractor", "./server.go")
	validate("github.com/viktorasm/gontractor/bar/foo", "bar/foo/server.go")
	validate("github.com/viktorasm/gontractor/bar/foo", "./bar/foo/server.go")

	// for absolute paths, should use the last src available
	validate("a/b/foo", "/disk/whatever/workspace/src/a/b/foo/server.go")
	validate("a/b/foo", "/disk/src/whatever/workspace/src/a/b/foo/server.go")
	validate("whatever/workspace/a/b/foo", "/disk/src/whatever/workspace/a/b/foo/server.go")

	// don't cough up on invalid srcs..
	validate("a/b/foosrc", "/disk/whatever/workspace/src/a/b/foosrc/server.go")
	validate("a/b/foosrc/sub", "/disk/whatever/workspace/src/a/b/foosrc/sub/server.go")
}
