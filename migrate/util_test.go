package migrate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const n = "test_for_build_plugin.go"
const s = "test_for_build_plugin.so"

func TestMain(m *testing.M) {
	MigrationTargetFolder = "."

	t := MigrationTargetFolder + "/" + n

	f, err := os.Create(t)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	if _, err := f.WriteString("package main"); err != nil {
		panic(err)
	}

	f.Sync()

	exitVal := m.Run()

	if err := os.Remove(t); err != nil {
		panic(err)
	}

	// Ensure the compiled files would be reset.
	if err := os.RemoveAll(MigrationTargetFolder + "/.plugins"); err != nil {
		panic(err)
	}

	os.Exit(exitVal)
}

func TestCheckPluginExists(t *testing.T) {
	path := MigrationTargetFolder + "/.plugins/" + s
	assert.False(t, checkPluginExists(s), "should return false when there is no plugin file exists")

	os.MkdirAll(MigrationTargetFolder+"/.plugins", 0755)
	os.Create(path)

	assert.True(t, checkPluginExists(s), "should return true since the plugin file exists")

	os.Remove(path)
}

func TestBuildPlugin(t *testing.T) {
	p, err := BuildPlugin(n)

	assert.NoError(t, err)
	assert.Equal(t, s, p, "should return builded file name which contains .so as ext")

	_, err = os.Stat(MigrationTargetFolder + "/.plugins/" + s)
	assert.Equal(t, false, os.IsNotExist(err), "the compiled plugin file should exists")
}
