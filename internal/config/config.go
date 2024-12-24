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
	URL     string `mapstructure:"url"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	KeyPath string `mapstructure:"key_path"`
}

type ApiConfig struct {
	BaseURL string `mapstructure:"base_url"`
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
