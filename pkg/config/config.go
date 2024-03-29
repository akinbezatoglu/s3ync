// https://github.com/cli/go-gh/tree/trunk/pkg/config
// Package config is a set of types for interacting with the configuration file.
package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/akinbezatoglu/s3ync/internal/yamlmap"
)

const (
	appData        = "AppData"
	s3yncConfigDir = "S3YNC_CONFIG_DIR"
	localAppData   = "LocalAppData"
	xdgConfigHome  = "XDG_CONFIG_HOME"
)

var (
	cfg     *Config
	once    sync.Once
	loadErr error
)

// Config is a in memory representation of the configuration file.
// It can be thought of as map where entries consist of a key that
// correspond to either a string value or a map value, allowing for
// multi-level maps.
type Config struct {
	entries *yamlmap.Map
	mu      sync.RWMutex
}

// Get a string value from a Config.
// The keys argument is a sequence of key values so that nested
// entries can be retrieved. A undefined string will be returned
// if trying to retrieve a key that corresponds to a map value.
// Returns "", KeyNotFoundError if any of the keys can not be found.
func (c *Config) Get(keys []string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := c.entries
	for _, key := range keys {
		var err error
		m, err = m.FindEntry(key)
		if err != nil {
			return "", &KeyNotFoundError{key}
		}
	}
	return m.Value, nil
}

// Keys enumerates a Config's keys.
// The keys argument is a sequence of key values so that nested
// map values can be have their keys enumerated.
// Returns nil, KeyNotFoundError if any of the keys can not be found.
func (c *Config) Keys(keys []string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := c.entries
	for _, key := range keys {
		var err error
		m, err = m.FindEntry(key)
		if err != nil {
			return nil, &KeyNotFoundError{key}
		}
	}
	return m.Keys(), nil
}

// Remove an entry from a Config.
// The keys argument is a sequence of key values so that nested
// entries can be removed. Removing an entry that has nested
// entries removes those also.
// Returns KeyNotFoundError if any of the keys can not be found.
func (c *Config) Remove(keys []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	m := c.entries
	for i := 0; i < len(keys)-1; i++ {
		var err error
		key := keys[i]
		m, err = m.FindEntry(key)
		if err != nil {
			return &KeyNotFoundError{key}
		}
	}
	err := m.RemoveEntry(keys[len(keys)-1])
	if err != nil {
		return &KeyNotFoundError{keys[len(keys)-1]}
	}
	return nil
}

// Set a string value in a Config.
// The keys argument is a sequence of key values so that nested
// entries can be set. If any of the keys do not exist they will
// be created. If the string value to be set is empty it will be
// represented as null not an empty string when written.
//
//	var c *Config
//	c.Set([]string{"key"}, "")
//	Write(c) // writes `key: ` not `key: ""`
func (c *Config) Set(keys []string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	m := c.entries
	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		entry, err := m.FindEntry(key)
		if err != nil {
			entry = yamlmap.MapValue()
			m.AddEntry(key, entry)
		}
		m = entry
	}
	val := yamlmap.StringValue(value)
	if value == "" {
		val = yamlmap.NullValue()
	}
	m.SetEntry(keys[len(keys)-1], val)
}

func (c *Config) deepCopy() *Config {
	return ReadFromString(c.entries.String())
}

// Read configuration file from the local file system and
// returns a Config. A copy of the fallback configuration will
// be returned when there are no configuration files to load.
// If there are no configuration files and no fallback configuration
// an empty configuration will be returned.
var Read = func(fallback *Config) (*Config, error) {
	once.Do(func() {
		cfg, loadErr = load(GeneralConfigFile(), fallback)
	})
	return cfg, loadErr
}

// ReadFromString takes a yaml string and returns a Config.
func ReadFromString(str string) *Config {
	m, _ := mapFromString(str)
	if m == nil {
		m = yamlmap.MapValue()
	}
	return &Config{entries: m}
}

// Write s3ync configuration files to the local file system.
// It will only write configuration file that have been modified
// since last being read.
func Write(c *Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.entries.IsModified() {
		err := writeFile(GeneralConfigFile(), []byte(c.entries.String()))
		if err != nil {
			return err
		}
		c.entries.SetUnmodified()
	}

	return nil
}

func load(generalFilePath string, fallback *Config) (*Config, error) {
	generalMap, err := mapFromFile(generalFilePath)
	if err != nil && !os.IsNotExist(err) {
		if errors.Is(err, yamlmap.ErrInvalidYaml) ||
			errors.Is(err, yamlmap.ErrInvalidFormat) {
			return nil, &InvalidConfigFileError{Path: generalFilePath, Err: err}
		}
		return nil, err
	}

	if generalMap == nil {
		generalMap = yamlmap.MapValue()
	}

	if generalMap.Empty() && fallback != nil {
		cfg := fallback.deepCopy()
		cfg.entries.SetModified()
		err := Write(cfg)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	// TODO: Check if the configuration files have the desired key values.
	// Suppose the user changed the configurations manually or for other reasons.

	return &Config{entries: generalMap}, nil
}

func GeneralConfigFile() string {
	return filepath.Join(ConfigDir(), "config.yml")
}

func mapFromFile(filename string) (*yamlmap.Map, error) {
	data, err := readFile(filename)
	if err != nil {
		return nil, err
	}
	return yamlmap.Unmarshal(data)
}

func mapFromString(str string) (*yamlmap.Map, error) {
	return yamlmap.Unmarshal([]byte(str))
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

func writeFile(filename string, data []byte) (writeErr error) {
	if writeErr = os.MkdirAll(filepath.Dir(filename), 0771); writeErr != nil {
		return
	}
	var file *os.File
	if file, writeErr = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600); writeErr != nil {
		return
	}
	defer func() {
		if err := file.Close(); writeErr == nil && err != nil {
			writeErr = err
		}
	}()
	_, writeErr = file.Write(data)
	return
}
