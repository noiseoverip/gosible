package internal

import (
	"ansiblego/internal/logging"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Logic related to variable handling

type GroupVariables map[string]map[string]interface{}

type HostVariables map[string]interface{}

func (hv HostVariables) Add(newVars map[string]interface{}) {
	for k, v := range newVars {
		hv[k] = v
	}
}

// load groups variables from provided groupByName vars directory structs
// for now only supports directory based groupByName vars
func LoadGroupVars(dir string) (groupVars GroupVariables, err error) {
	groupVars = make(map[string]map[string]interface{})
	root := filepath.Join(dir, "group_vars")
	// Exist if doesn't exist
	if _, err := os.Stat(root); os.IsNotExist(err) {
		logging.Info("no group_vars dir found, skipping\n")
		return nil, nil
	}
	// Assemble a list of tier variable files
	err = filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			logging.Debug("Loading group vars from: %s\n", filePath)
			group := path.Base(path.Dir(filePath))

			var variables = make(map[string]interface{})

			varsFileReader, err := os.Open(filePath)
			if err != nil {
				return nil
			}

			contents, err := ioutil.ReadAll(varsFileReader)
			if err != nil {
				return nil
			}

			if err := yaml.Unmarshal(contents, variables); err != nil {
				return fmt.Errorf("failed to load variables from %s", filePath)
			}

			if _, ok := groupVars[group]; !ok {
				groupVars[group] = make(map[string]interface{})
			}

			for k, v := range variables {
				groupVars[group][k] = v
			}

		}
		return nil
	})
	return
}
