package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

type Config struct {
	SSH SSHConfig `mapstructure:"ssh"`
	Api ApiConfig `mapstructure:"api"`
}

type SSHConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ApiConfig struct {
	BaseUrl string `mapstructure:"base_url"`
}

func New() Config {
	viper.SetConfigFile(os.Getenv("CONFIG_FILE"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return config
}