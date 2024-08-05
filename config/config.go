package config

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	HTTP   HTTP   `yaml:"http"`
	DB     DB     `yaml:"db"`
	Logger Logger `yaml:"logger"`
}

type HTTP struct {
	Addr string `yaml:"addr"`
}

type DB struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"-"`
}

type Logger struct {
	Production        bool `yaml:"production"`
	DisableStacktrace bool `yaml:"disable_stacktrace"`
}

func Must() *Config {
	cfg := new(Config)

	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(data, cfg); err != nil {
		panic(err)
	}

	if err = godotenv.Load(); err != nil {
		panic(err)
	}

	cfg.DB.DSN = os.Getenv("DSN")

	return cfg
}
