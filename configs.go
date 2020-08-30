package main

import "time"

type ApplicationConfiguration struct {
	Development       bool   `yaml:"development"`
	Database          string `yaml:"database"`
	ListenOnPort      string `yaml:"listenOnPort"`
	TemplateDirectory string `yaml:"templateDirectory"`
	WriteTimeout      time.Duration    `yaml:"writeTimeout"`
	ReadTimeout        time.Duration   `yaml:"readTimeout"`
}

type Properties struct {
	Values map[string]string `yaml:"properties"`
}

type String string
type Double float64
type Integer int64
type SmallInt int32
type Boolean bool

const (
	DefaultConfigurationPath String = "config/local.yaml"
)
