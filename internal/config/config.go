package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

var (
	config *Config
)

const (
	// Environment constants
	development = "development"
	staging     = "staging"
	production  = "production"
	// SERVICE constant for service configuration load
	SERVICE = "multiclient-websocket"
)

type option struct {
	configFile string
}

// Init ...
func Init(opts ...Option) error {
	opt := &option{
		configFile: getDefaultConfigFile(),
	}
	for _, optFunc := range opts {
		optFunc(opt)
	}

	out, err := ioutil.ReadFile(opt.configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(out, &config)
}

// Option ...
type Option func(*option)

// WithConfigFile ...
func WithConfigFile(file string) Option {
	return func(opt *option) {
		opt.configFile = file
	}
}

func getDefaultConfigFile() string {
	env := "development"
	configPath := filepath.Join(getRootDir(), "files/etc/"+SERVICE)
	// For Kubernetes namespaces
	// namespace, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	// env = string(namespace)
	//

	switch env {
	case development:
		// Docker Env
		if os.Getenv("GOPATH") == "" {
			configPath = "./" + SERVICE + ".development.yaml"
		} else {
			configPath = filepath.Join(configPath, "/"+SERVICE+".development.yaml")
		}
	case staging:
		configPath = filepath.Join(configPath, "/"+SERVICE+".staging.yaml")
	case production:
		configPath = filepath.Join(configPath, "/"+SERVICE+".production.yaml")
	}
	return configPath
}

// Get ...
func Get() *Config {
	return config
}

func getRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Join(filepath.Dir(d), "..")
}
