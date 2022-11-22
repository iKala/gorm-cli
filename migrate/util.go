package migrate

import (
	"fmt"
	"os"
	"os/exec"
	"plugin"
	"strings"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func checkPluginExists(pluginName string) bool {
	_, err := os.Stat(MigrationTargetFolder + "/.plugins/" + pluginName)
	return !os.IsNotExist(err)
}

func BuildPlugin(fileName string) (string, error) {
	goFileName := strings.Replace(fileName, ".so", ".go", -1)
	pluginName := strings.Replace(fileName, ".go", ".so", -1)

	// Skip build plugin when exists.
	if checkPluginExists(pluginName) {
		return pluginName, nil
	}

	cmd := exec.Command(
		"go",
		"build",
		"-buildmode=plugin",
		"-o",
		MigrationTargetFolder+"/.plugins/"+pluginName,
		MigrationTargetFolder+"/"+goFileName,
	)

	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return pluginName, nil
}

func getPlugin(pluginName string, symbolTarget string) (plugin.Symbol, error) {
	plug, err := plugin.Open(MigrationTargetFolder + "/.plugins/" + pluginName)
	if err != nil {
		return nil, fmt.Errorf("%v (%v)", "Open plugin filed.", err.Error())
	}

	s, err := plug.Lookup(symbolTarget)
	if err != nil {
		return nil, fmt.Errorf("%v (%v)", symbolTarget+"wrong format - missing "+symbolTarget+" declaration.", err.Error())
	}

	return s, nil
}

func getMigration(pluginName string) (Migration, error) {
	s, err := getPlugin(pluginName, "Migration")
	if err != nil {
		return nil, err
	}

	var migration Migration
	migration, ok := s.(Migration)
	if !ok {
		return nil, fmt.Errorf("%v (%v)", "Unexpected type from module symbol.", err.Error())
	}
	return migration, nil
}

func RemovePlugin(goFileName string) error {
	pluginName := strings.Replace(goFileName, ".go", ".so", -1)

	if !checkPluginExists(pluginName) {
		return nil
	}

	return os.Remove(MigrationTargetFolder + "/.plugins/" + pluginName)
}

// MigrationTargetFolder is the migration folder target.
var MigrationTargetFolder string
