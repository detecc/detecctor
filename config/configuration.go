package config

import (
	"flag"
	"fmt"
	"github.com/kkyr/fig"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"path/filepath"
	cache2 "github.com/detecc/detecctor/cache"
)

// GetFlags get the program flags and store them in the cache
func GetFlags() {
	var (
		memory = cache2.Memory()
	)
	var workingDirectory, _ = os.Getwd()
	var configurationPath = fmt.Sprintf("%s", workingDirectory)

	// Get the paths from arguments
	configurationFileFormatFlag := flag.String("config-format", "yaml", "Format of the configuration files (YAML, JSON or TOML)")
	configurationFilePathFlag := flag.String("config-file", fmt.Sprintf("%s/config.%s", configurationPath, *configurationFileFormatFlag), "Path of the configuration file")
	flag.Parse()

	// Put the configuration file path into the cache
	memory.Set("configurationFilePath", *configurationFilePathFlag, cache.NoExpiration)
}

// GetServerConfiguration get the configuration from the configuration file and store the configuration in the cache
func GetServerConfiguration() *Configuration {
	var (
		config                Configuration
		err                   error
		configurationFilePath string
		memory                = cache2.Memory()
		isFound               bool
		cachedConfiguration   interface{}
	)

	cachedConfiguration, isFound = memory.Get("configuration")
	if isFound {
		return cachedConfiguration.(*Configuration)
	}

	configurationPath, isFound := memory.Get("configurationFilePath")
	if isFound {
		configurationFilePath = configurationPath.(string)
	} else {
		log.Fatal("No configuration file path found!")
	}

	err = fig.Load(&config,
		fig.File(filepath.Base(configurationFilePath)),
		fig.Dirs(filepath.Dir(configurationFilePath)),
	)
	if err != nil {
		panic(err)
	}

	memory.Set("configuration", &config, cache.NoExpiration)

	return &config
}
