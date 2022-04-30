package config

import (
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DestFileFormat  string
	DestDirectory   string
	SourceDirectory string
	NoRecurse       bool
	Verbose         bool
	Workers         int
}

func (c *Config) PrintConfig() {
	log.Infof("Running with config: ")
	log.Infof(" * DestFileFormat = %v", c.DestFileFormat)
	log.Infof(" * DestDirectory = %v", c.DestDirectory)
	log.Infof(" * SourceDirectory = %v", c.SourceDirectory)
	log.Infof(" * NoRecurse = %v", c.NoRecurse)
	log.Infof(" * Verbose = %v", c.Verbose)
}
