package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DbHost    string `mapstructure:"DB_HOST"`
	DbPort    int    `mapstructure:"DB_PORT"`
	DbUser    string `mapstructure:"DB_USER"`
	DbPass    string `mapstructure:"DB_PASSWORD"`
	DbName    string `mapstructure:"DB_NAME"`
	AppPort   string `mapstructure:"APP_PORT"`
	Token	  string `mapstructure:"TOKEN"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
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
