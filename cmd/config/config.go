package config

import (
	"github.com/mitchellh/go-homedir"
	"log"
)

func (c *Config) ConfigDir() string {
	c.configDirInit.Do(func() {
		var err error
		if c.configDir, err = homedir.Expand("~/.orca/"); err != nil {
			log.Fatal("failed to get config dir due to", err)
		}
	})
	return c.configDir
}
