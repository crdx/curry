package config

import (
	"log"
	"os"
	"path/filepath"

	"crdx.org/col"

	"github.com/BurntSushi/toml"
)

const toolName = "org.crdx/curry"

type Config struct {
	APIKey string `toml:"api_key"`
}

func Load() Config {
	configHome, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(col.Red("Error: unable to determine config directory: " + err.Error()))
	}

	configPath := filepath.Join(configHome, toolName, "config.toml")

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal(col.Red("Error: unable to read config: " + err.Error()))
	}

	if config.APIKey == "" {
		log.Fatal(col.Red("Error: api_key is not set in " + configPath))
	}

	return config
}
