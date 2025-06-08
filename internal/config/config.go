package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTPServer HTTPServerConfig `yaml:"http_server"`
	Postgres   PostgresConfig   `yaml:"postgres"`
	Ollama     OllamaConfig     `yaml:"ollama"`
}

type HTTPServerConfig struct {
	Host           string        `yaml:"host"`
	Port           string        `yaml:"port"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes int           `yaml:"max_header_bytes"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type OllamaConfig struct {
	Host string `yaml:"host"`
}

func MustLoad(configPath string) *Config {
	if configPath == "" {
		panic("config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	file, err := os.ReadFile(configPath)
	if err != nil {
		panic("failed to read config file: " + err.Error())
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		panic("failed to parse config file: " + err.Error())
	}

	return &cfg
}
