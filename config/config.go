package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

const (
	configPathKey     = "CONFIG_PATH"
	defaultConfigPath = ".env"
)

type Config struct {
	DB  DBConfig  `yaml:"db" env-prefix:"DB_"`
	JWT JWTConfig `yaml:"jwt" env-prefix:"JWT_"`
}

type DBConfig struct {
	Port     int    `yaml:"port" env:"PORT" env-default:"3306"`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"NAME" env-default:"clean"`
	User     string `yaml:"user" env:"USER" env-default:"root"`
	Password string `yaml:"password" env:"PASSWORD"`
}

type JWTConfig struct {
	Secret string `yaml:"secret" env:"SECRET"`
}

type FilePath string

func ProvideConfig(configPath FilePath) (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig(string(configPath), &cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("Application configuration initialize failed, %s has a problem", configPath)
		return cfg, err
	}
	return cfg, nil
}

func GetConfigPath() FilePath {
	env, present := os.LookupEnv(configPathKey)
	if present {
		return FilePath(env)
	}
	return defaultConfigPath
}
