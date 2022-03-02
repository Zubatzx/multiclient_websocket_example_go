package config

type (
	// Config ...
	Config struct {
		ProjectName string       `yaml:"projectName"`
		ServiceName string       `yaml:"serviceName"`
		Env         string       `yaml:"env"`
		Server      ServerConfig `yaml:"server"`
	}

	// ServerConfig ...
	ServerConfig struct {
		Port string `yaml:"port"`
	}
)
