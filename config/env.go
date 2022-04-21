package config

import (
	"github.com/caarlos0/env"
	"github.com/go-playground/log"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AppName   string `env:"APP_NAME" envDefault:"podcastapi"`
	AppEnv    string `env:"APP_ENV" envDefault:"local"`
	AppDebug  string `env:"APP_DEBUG" envDefault:"true"`
	AppPort   string `env:"APP_PORT" envDefault:"8080"`
	AppDomain string `env:"APP_DOMAIN" envDefault:"0.0.0.0"`

	DBName string `env:"DB_NAME" envDefault:"podscraper"`
	DBPass string `env:"DB_PASS" envDefault:"Underwood42"`
	DBUser string `env:"DB_USER" envDefault:"root"`
	DBType string `env:"DB_TYPE" envDefault:"mysql"`
	DBHost string `env:"DB_HOST" envDefault:"172.17.0.1"`
	DBPort string `env:"DB_PORT" envDefault:"3306"`

	RedisHost string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPass string `env:"REDIS_PASS" envDefault:""`

	JWTSecret     string `env:"JWT_SECRET" envDefault:"qBPXnbcuQyauqlhTpJQjgAnmauKiZUgrhdu7eQhuNXfr6"`
	MailGunApiKey string `env:"MAILGUN_API_KEY" envDefault:""`
}

var cfg *Config

// Parse parses, validates and then returns the application
// configuration based on ENV variables
func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Warn("File .env not found, reading configuration from ENV")
	}

	log.Info("Parsing ENV vars...")
	defer log.Info("Done Parsing ENV vars")

	cfg = &Config{}

	if err := env.Parse(cfg); err != nil {
		log.WithFields(log.F("error", err)).Warn("Errors Parsing Configuration")
	}

	return
}

func GetConfig() *Config {
	return cfg
}
