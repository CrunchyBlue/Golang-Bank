package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	AccessTokenDuration     time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	AccessTokenSymmetricKey string        `mapstructure:"ACCESS_TOKEN_SYMMETRIC_KEY"`
	DBDriver                string        `mapstructure:"DB_DRIVER"`
	DBSource                string        `mapstructure:"DB_SOURCE"`
	RefreshTokenDuration    time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	ServerAddress           string        `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
