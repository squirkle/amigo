package amigo

import (
	"os"

	"github.com/pelletier/go-toml"
)

const EnvMapKey = "envmap"

type Config struct {
	env      map[string]string
	confFile *toml.TomlTree
}

// Associate a given config file key with an environment var
func (c *Config) Env(confKey string, envKey string) {
	if val := os.Getenv(envKey); val != "" {
		c.env[confKey] = val
	}
}

func (c *Config) Get(key string) interface{} {
	if val, ok := c.env[key]; ok {
		return val
	}
	return c.confFile.Get(key)
}

// Return a new configuration object for use by library consumers
func New(filepath string) (*Config, error) {
	file, err := toml.LoadFile(filepath)

	if err != nil {
		return nil, err
	}

	c := &Config{}
	c.env = make(map[string]string)
	c.confFile = file

	// if a envmap table is defined, associate the specified keys with the env
	// vars defined there
	if file.Has(EnvMapKey) {
		envmap := file.Get(EnvMapKey).(*toml.TomlTree)
		for _, confKey := range envmap.Keys() {
			envKey := envmap.Get(confKey).(string)
			c.Env(confKey, envKey)
		}
	}
	return c, nil
}
