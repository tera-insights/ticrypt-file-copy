package config

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
)

const CONFIG_FILE = "/etc/ticp/ticp.conf"

type server struct {
	AllowedHosts []string `toml:"allowed_hosts"`
	Port         string   `toml:"port"`
}

type storage struct {
	Db   string `toml:"db"`
	Path string `toml:"path"`
}

type copyConfig struct {
	ChunkSize int `toml:"chunk_size"`
}

type Config struct {
	Server  server     `toml:"server"`
	Storage storage    `toml:"storage"`
	Copy    copyConfig `toml:"copy"`
}

func FetchConfig() *Config {

	returnDefaultConfig := func(err error) *Config {
		fmt.Printf("[Warning]Error reading config file: %v: Using defaults\n", err)

		// Return the default config
		return &Config{
			Server: server{
				AllowedHosts: []string{"localhost"},
				Port:         "4242",
			},
			Storage: storage{
				Db:   "ticp.db",
				Path: "/var/lib/ticp",
			},
			Copy: copyConfig{
				ChunkSize: 4,
			},
		}
	}

	// Read the config file
	content, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		return returnDefaultConfig(err)
	}

	// Parse the config file
	var config Config
	err = toml.Unmarshal(content, &config)
	if err != nil {
		return returnDefaultConfig(err)
	}
	// Return the config
	return &config
}
