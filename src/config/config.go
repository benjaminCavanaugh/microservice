package config

var fileName string

type Config struct {
	serviceConfig ServiceConfig
}

func NewConfig() Config {
	var c Config;
	return c;
}

func (c Config) ParseFromFile() {
	// TODO: Write a function for parsing a configuration YAML
}

func (c Config) GetConfig() ServiceConfig {
	return c.serviceConfig
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
	ReadTimeout  string
	WriteTimeout string
}

func NewHtmlServerConfig() HtmlServerConfig {
	var h HtmlServerConfig;
	return h;
}