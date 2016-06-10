package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

var configDirectoryPath string

// ServiceConfig contains configuration information for a service
type ServiceConfig struct {
	Token string
}

// Config contains configuration information
type Config struct {
	DatabasePath string                    `yaml:"databasePath"`
	Services     map[string]*ServiceConfig `yaml:"services"`
}

// GetService returns the configuration information for a service
func (config *Config) GetService(name string) *ServiceConfig {
	if config.Services == nil {
		config.Services = make(map[string]*ServiceConfig)
	}

	service := config.Services[name]
	if service == nil {
		service = &ServiceConfig{}
		config.Services[name] = service
	}
	return service
}

// ReadConfig reads the configuration information
func ReadConfig() (*Config, error) {
	file := configFilePath()

	var config Config
	if _, err := os.Stat(file); err == nil {
		// Read and unmarshal file only if it exists
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(f, &config)
		if err != nil {
			return nil, err
		}
	}

	// Set default database path
	if config.DatabasePath == "" {
		config.DatabasePath = path.Join(configDirectoryPath, fmt.Sprintf("%s.db", ProgramName))
	}
	return &config, nil
}

// WriteConfig writes the configuration information
func (config *Config) WriteConfig() error {
	err := os.MkdirAll(configDirectoryPath, 0700)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFilePath(), data, 0600)
}

func configFilePath() string {
	return path.Join(configDirectoryPath, fmt.Sprintf("%s.yaml", ProgramName))
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		// TODO Die

	} else {
		configDirectoryPath = path.Join(home, ".config", ProgramName)
	}
}
