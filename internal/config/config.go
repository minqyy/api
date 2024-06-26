package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/minqyy/api/pkg/slice"
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
	Env      string   `yaml:"env" env-required:"true"`
	Token    Token    `yaml:"token"`
	Server   Server   `yaml:"server" env-required:"true"`
	Hasher   Hasher   `yaml:"hasher"`
	Postgres Postgres `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
}

type Token struct {
	Access  TokenAccess  `yaml:"access"`
	Refresh TokenRefresh `yaml:"refresh"`
}

type TokenAccess struct {
	Secret string        `yaml:"secret"`
	TTL    time.Duration `yaml:"ttl"`
}

type TokenRefresh struct {
	TTL time.Duration `yaml:"ttl"`
}

type Server struct {
	Address     string        `yaml:"address" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Hasher struct {
	Salt string `yaml:"salt"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	ModeSSL  string `yaml:"sslmode"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatalf("Missed CONFIG_PATH environment parameter\n")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at: %s\n", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Cannot read config at %s: %s\n", configPath, err)
	}

	if !slice.Contains([]string{EnvLocal, EnvDevelopment, EnvProduction}, config.Env) {
		log.Fatalf("Unknown environment parameter. Use: `local`, `dev` or `prod` values\n")
	}

	return &config
}
