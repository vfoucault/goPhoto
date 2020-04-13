package config

import (
	"flag"
	"github.com/vfoucault/goPhoto/logger"
	"os"
	"runtime"
)

type Config struct {
	DestFileFormat  string
	DestDirectory   string
	SourceDirectory string
	NoRecurse       bool
	Verbose         bool
	Workers         int
}

func (c *Config) Init() *Config {
	flag.StringVar(&c.DestFileFormat, "format", "2006/2006-01-02", "Destination File format")
	flag.StringVar(&c.SourceDirectory, "source", "./", "the source directory")
	flag.BoolVar(&c.NoRecurse, "no-recurse", false, "Don't search recursively")
	flag.StringVar(&c.DestDirectory, "destination", "", "the destination directory")
	flag.BoolVar(&c.Verbose, "verbose", false, "be verbose")
	flag.IntVar(&c.Workers, "num-workers", runtime.NumCPU(), "number of workers. Default to runtime.NumCPU()")

	flag.Parse()

	if c.DestDirectory == "" {
		logger.Errorf("Destination direction is empty")
		os.Exit(1)
	}

	return c
}

func (c *Config) PrintConfig() {
	logger.Infof("Running with config: ")
	logger.Infof(" * DestFileFormat = %v", c.DestFileFormat)
	logger.Infof(" * DestDirectory = %v", c.DestDirectory)
	logger.Infof(" * SourceDirectory = %v", c.SourceDirectory)
	logger.Infof(" * NoRecurse = %v", c.NoRecurse)
	logger.Infof(" * Verbose = %v", c.Verbose)
}
