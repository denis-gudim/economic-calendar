package api

import (
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
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
		return xerrors.Errorf("error reading config file: %w", err)
	}

	return viper.Unmarshal(cnf)
}
