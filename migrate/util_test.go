package migrate

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const n = "test_for_build_plugin.go"

func TestMain(m *testing.M) {
	MigrationTargetFolder = "./"

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

func TestBuildPlugin(t *testing.T) {
	p, err := buildPlugin(n)

	assert.NoError(t, err)
	assert.Equal(t, "test_for_build_plugin.so", p, "should return builded file name which contains .so as ext")

	_, err = os.Stat(MigrationTargetFolder + "/.plugins/test_for_build_plugin.so")
	assert.Equal(t, false, os.IsNotExist(err), "the compiled plugin file should exists")
}
