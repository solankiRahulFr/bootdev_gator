// reading and writing the json file
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// const configFileName = "/bootdev/bootdev_gator/.gatorconfig.json"
const configFileName = "/.gatorconfig.json"

// Config struct that represents the JSON file structure
type Config struct {
	Dburl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}

func Read() (Config, error) {
	// read the json file and return the config struct
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return err
	}
	return nil
}

func (c *Config) SetUserName(name string) error {
	c.CurrentUser = name
	return write(*c)
}
