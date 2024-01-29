package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

const (
	appData        = "AppData"
	s3yncConfigDir = "S3YNC_CONFIG_DIR"
	localAppData   = "LocalAppData"
	xdgConfigHome  = "XDG_CONFIG_HOME"
)

var path = filepath.Join(ConfigDir(), "config.yml")

type Syncs struct {
	// Key: filesytem path, Value: [aws-profile, bucket-name]
	All map[string][]string
}

func GetAllSyncList() *Syncs {
	data, _ := readFile(path)

	var syncs Syncs
	syncs.All = make(map[string][]string)

	// Create a map to hold the parsed YAML data
	var yamlData map[string]interface{}
	// Unmarshal the YAML text into the map
	err := yaml.Unmarshal(data, &yamlData)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML data: %v", err)
	}

	m, err := extractNestedValue(yamlData, []string{"s3", "profiles"}...)
	if err != nil {
		fmt.Println(err)
	}
	for v := range m {
		m, err := extractNestedValue(yamlData, []string{"s3", "profiles", v, "syncs"}...)
		if err != nil {
			fmt.Println(err)
		}
		for vv := range m {
			m, err := extractNestedValue(yamlData, []string{"s3", "profiles", v, "syncs", vv}...)
			if err != nil {
				fmt.Println(err)
			}
			syncs.All[m["local"].(string)] = []string{v, m["bucket"].(string)}
		}
	}
	return &syncs
}

func extractNestedValue(data map[string]interface{}, keys ...string) (map[string]interface{}, error) {
	var currentMap map[string]interface{} = data

	for _, key := range keys {
		value, found := currentMap[key]
		if !found {
			return nil, fmt.Errorf("key '%s' not found", key)
		}

		nextMap := value.(map[string]interface{})
		currentMap = nextMap
	}

	return currentMap, nil
}

func readFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Config path precedence: S3YNC_CONFIG_DIR, XDG_CONFIG_HOME, AppData (windows only), HOME.
func ConfigDir() string {
	var path string
	if a := os.Getenv(s3yncConfigDir); a != "" {
		path = a
	} else if b := os.Getenv(xdgConfigHome); b != "" {
		path = filepath.Join(b, "s3ync")
	} else if c := os.Getenv(appData); runtime.GOOS == "windows" && c != "" {
		path = filepath.Join(c, "s3ync")
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, ".config", "s3ync")
	}
	return path
}
