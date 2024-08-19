package main

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

type config struct {
	Server  server     `toml:"server"`
	Storage storage    `toml:"storage"`
	Copy    copyConfig `toml:"copy"`
}

func fetchConfig() *config {
	// Read the config file
	content, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		fmt.Printf("[Warning]Error reading config file: %v: Using defaults\n", err)

		// Return the default config
		return &config{
			Server: server{
				AllowedHosts: []string{"localhost"},
				Port:         "4242",
			},
			Storage: storage{
				Db:   "ticp.db",
				Path: "/var/lib/ticp",
			},
			Copy: copyConfig{
				ChunkSize: 4 * 1024 * 1024,
			},
		}
	}

	// Parse the config file
	var config *config
	toml.Unmarshal(content, config)
	// Return the config
	return config
}
