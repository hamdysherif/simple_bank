package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config store all configuration
type Config struct {
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	DBSourceTest  string        `mapstructure:"DB_SOURCE_TEST"`
	ServerAddress string        `mapstructure:"SERVER_ADDRESS"`
	SemmetricKey  string        `mapstructure:"SYMMETRIC_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

// LoadConfig to return all configuration
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		return
	}

	err = viper.Unmarshal(&config)
	return
}
