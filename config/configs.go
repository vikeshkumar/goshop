package config

import (
	"time"
)

// Configuration structure for the application
type Configuration struct {
	ServerConfiguration   ServerConfiguration   `yaml:"serverConfiguration"`
	DBConfig              DBConfig              `yaml:"dbConfig"`
	TemplateConfiguration TemplateConfiguration `yaml:"templateConfiguration"`
}

// ServerConfiguration is the
type ServerConfiguration struct {
	Development  bool          `yaml:"development"`
	ListenOnPort string        `yaml:"listenOnPort"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
}

// TemplateConfiguration is the configuration for the templates
type DBConfig struct {
	URL string `yaml:"url"`
}

// TemplateConfiguration is the configuration for the templates
type TemplateConfiguration struct {
	ParseOnce         bool              `yaml:"parseOnce"`
	TemplateDirectory string            `yaml:"templateDirectory"`
	Favicon           string            `yaml:"favicon"`
	Templates         map[string]string `yaml:"templates"`
}

// String as string for use in constants
type String string

// constants for the application will be used as default values
const (
	DefaultConfigurationPath String = "config/local.yaml"
)
