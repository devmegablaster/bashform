package config

import (
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

type Config struct {
	SSH      SSHConfig      `mapstructure:"ssh"`
	Database DatabaseConfig `mapstructure:"database"`
	Crypto   CryptoConfig   `mapstructure:"crypto"`
}

type SSHConfig struct {
	URL     string `mapstructure:"url"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	KeyPath string `mapstructure:"key_path"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type CryptoConfig struct {
	AESKey string `mapstructure:"aes_key"`
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

	config.loadEnv()

	slog.Info("✅ Config loaded")

	return config
}

func (c *Config) loadEnv() {
	// database secrets
	c.Database.Host = os.ExpandEnv(c.Database.Host)
	c.Database.Port = os.ExpandEnv(c.Database.Port)
	c.Database.User = os.ExpandEnv(c.Database.User)
	c.Database.Password = os.ExpandEnv(c.Database.Password)
	c.Database.Name = os.ExpandEnv(c.Database.Name)

}
