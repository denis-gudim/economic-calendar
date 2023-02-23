package api

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		ConnectionString string `mapstructure:"DB_CONSTR"`
	} `mapstructure:",squash"`
}

func (cnf *Config) Load() error {

	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("couldn't read configuration file: %w", err)
	}

	return viper.Unmarshal(cnf)
}
