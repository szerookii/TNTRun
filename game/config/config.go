package config

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/pelletier/go-toml"
	"os"
)

// Config ...
type Config struct {
	Enabled bool       `json:"enabled"`
	Lobby   mgl64.Vec3 `json:"lobby"`
}

// GetConfig ...
func GetConfig() (*Config, error) {
	if _, err := os.Stat("tntrun.toml"); os.IsNotExist(err) {
		config := &Config{}
		config.Enabled = false
		config.Lobby = mgl64.Vec3{}

		bytes, _ := toml.Marshal(config)

		if err = os.WriteFile("tntrun.toml", bytes, 0777); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile("tntrun.toml")

	if err != nil {
		return nil, err
	}

	config := &Config{}

	_ = toml.Unmarshal(data, config)

	return config, nil
}

// UpdateConfig ...
func UpdateConfig(enabled bool, lobby mgl64.Vec3) error {
	config := &Config{}
	config.Enabled = enabled
	config.Lobby = lobby

	bytes, _ := toml.Marshal(config)

	if err := os.WriteFile("tntrun.toml", bytes, 0777); err != nil {
		return err
	}

	return nil
}
