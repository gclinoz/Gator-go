package config

import (
	"os"
	"encoding/json"
	"path/filepath"
)

const configFileNmae = ".gatorconfig.json"

type Config struct {
	DBURL		string	`json:"db_url"`
	Username	string	`json:"current_user_name"`
}

func (c *Config) SetUser(name string) error {
	c.Username = name
	return write(*c)
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var data Config 
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return Config{}, err
	}
	return data, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullpath := filepath.Join(home, configFileNmae)
	return fullpath, nil
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0600)
	return err
}
