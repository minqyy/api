package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const (
	EnvLocal       = "local"
	EnvDevelopment = "dev"
	EnvProduction  = "prod"
)

type Config struct {
	Env    string `yaml:"env" env-required:"true"`
	Server Server `yaml:"server" env-required:"true"`
}

//type PostgresDB struct {
//	Host     string `yaml:"host"`
//	Port     string `yaml:"port"`
//	User     string `yaml:"user"`
//	Name     string `yaml:"name"`
//	Password string `yaml:"password"`
//	ModeSSL  string `yaml:"sslmode"`
//}

//type Redis struct {
//	Host     string `yaml:"host"`
//	Username string `yaml:"username"`
//	Password string `yaml:"password"`
//}

type Server struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// MustLoad loads config to a new Config instance and return it.
func MustLoad() *Config {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatalf("missed CONFIG_PATH environment parameter")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist at: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config at %s: %s", configPath, err)
	}

	return &config
}
