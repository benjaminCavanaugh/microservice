package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

func HelloWorld() {
	fmt.Println("Hello World!");
}

type Config struct {
	serviceConfig ServiceConfig
}

func NewConfig() Config {
	var c Config;
	return c;
}

func ParseFromFile(filename string) (*Config, error) {
    buf, err := ioutil.ReadFile(filename)

	if err != nil {
        return nil, err
    }

    c := &Config{}
    err = yaml.Unmarshal(buf, c)

	if err != nil {
        return nil, fmt.Errorf("in file %q: %w", filename, err)
    }

    return c, err
}

func (c Config) GetServiceConfig() (ServiceConfig, error) {
	// if c.serviceConfig == nil {
	// 	return nil, error.Error()
	// }

	return c.serviceConfig, nil
}

type ServiceConfig struct {
	htmlServerConfig HtmlServerConfig
}

func NewServiceConfig() ServiceConfig {
	var s ServiceConfig;
	return s;
}

func (s ServiceConfig) GetHtmlServerConfig() HtmlServerConfig {
	return s.htmlServerConfig;
}

type HtmlServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewHtmlServerConfig() HtmlServerConfig {
	var h HtmlServerConfig;
	return h;
}