package config

import (
	"errors"
	"fmt"
	"os"

	"flag"

	"github.com/PeterMue/bitbucket-webhook/handler"
	"gopkg.in/yaml.v3"
)

type Hooks struct {
	EventKey string   `yaml:"event"`
	Command  string   `yaml:"command"`
	Args     []string `yaml:"args"`
	Async    bool     `yaml:"async"`
}

type Config struct {
	Secret string  `yaml:"secret"`
	Listen string  `yaml:"listen"`
	Hooks  []Hooks `yaml:"hooks"`
}

func (cfg Config) Handler(eventKey string) (*handler.Handler, error) {
	for _, hook := range cfg.Hooks {
		if hook.EventKey == eventKey {
			return handler.New(hook.EventKey, hook.Command, hook.Args, hook.Async), nil
		}
	}
	return nil, errors.New("no handler found")
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ParseFlags(args []string) (*Config, error) {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	// Define flags for command line or environment
	var (
		file   = flags.String("config", "config.yaml", "Path to config file")
		listen = flags.String("listen", "", "Server listening address")
		secret = flags.String("secret", "", "Webhook secret to verify signature")
	)

	// Prase flags
	if err := flags.Parse(args[1:]); err != nil {
		return nil, err
	}

	// Validate file location
	if err := validateFile(*file); err != nil {
		return nil, err
	}

	// Parse config file
	config, err := NewConfig(*file)
	if err != nil {
		return nil, err
	}

	// Override config settings if passed via command line or environment
	if *listen != "" {
		config.Listen = *listen
	}
	if *secret != "" {
		config.Secret = *secret
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *Config) error {
	missing := func(name string) error {
		return fmt.Errorf("Missing configuration property: %s", name)
	}
	if config.Listen == "" {
		return missing("listen")
	}
	if config.Secret == "" {
		return missing("secret")
	}
	if len(config.Hooks) == 0 {
		return missing("hooks")
	}
	for i, hook := range config.Hooks {
		if hook.Command == "" {
			return missing(fmt.Sprintf("hooks[%d].command", i))
		}
		if hook.EventKey == "" {
			return missing(fmt.Sprintf("hooks[%d].event", i))
		}
	}
	return nil
}

func validateFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("Config file '%s' is a directory.", path)
	}
	return nil
}
