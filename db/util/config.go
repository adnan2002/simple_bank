package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbHost    string `mapstructure:"DB_HOST"`
	DbPort    int    `mapstructure:"DB_PORT"`
	DbUser    string `mapstructure:"DB_USER"`
	DbPass    string `mapstructure:"DB_PASSWORD"`
	DbSslMode string `mapstructure:"DB_SSL_MODE"`
	DbName    string `mapstructure:"DB_NAME"`
	AppPort   string `mapstructure:"APP_PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
