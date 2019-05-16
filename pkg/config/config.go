package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Config represents the configuration for wuclient.
type Config struct {
	AuthserverURL string `json:"authserver_url" long:"authserver-url"`
	WuserverURL   string `json:"wuserver_url" long:"wuserver-url"`
}

// ParseConfig parses and returns config struct.
func ParseConfig() *Config {
	c := &Config{
		AuthserverURL: "localhost:10000",
		WuserverURL:   "localhost:10000",
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		panic("failed to get user home")
	}

	wrapupsDir := path.Join(userHome, ".wrapups")
	configPath := path.Join(wrapupsDir, "config")
	if _, err = os.Stat(configPath); err == nil {
		b, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic("failed to read config")
		}
		if err := json.Unmarshal(b, &c); err != nil {
			panic("failed to parse config file")
		}
	}

	return c
}
