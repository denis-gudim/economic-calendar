package loader

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/xerrors"
)

type Config struct {
	DB struct {
		ConnectionString string `mapstructure:"DB_CONSTR"`
	} `mapstructure:",squash"`
	Loading struct {
		DefaultLanguageId int       `mapstructure:"LOADING_BATCHSIZE"`
		RetryCount        int       `mapstructure:"LOADING_DEFAULTLANG"`
		BatchSize         int       `mapstructure:"LOADING_RETRYCOUNT"`
		FromTime          time.Time `mapstructure:"LOADING_FROMTIME"`
		ToDays            int       `mapstructure:"LOADING_TODAYS"`
	} `mapstructure:",squash"`
	Logging struct {
		Level log.Level `mapstructure:"LOG_LEVEL"`
	} `mapstructure:",squash"`
	Scheduler struct {
		HistoryExpression string `mapstructure:"SCHEDULER_HISTEXPR"`
		RefreshExpression string `mapstructure:"SCHEDULER_REFREXPR"`
	} `mapstructure:",squash"`
}

func (cnf *Config) Load() error {

	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return xerrors.Errorf("Error reading config file: %w", err)
	}

	return viper.Unmarshal(cnf, func(m *mapstructure.DecoderConfig) {
		m.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			func(
				f reflect.Type,
				t reflect.Type,
				data interface{}) (interface{}, error) {
				if f.Kind() != reflect.String {
					return data, nil
				}
				if t != reflect.TypeOf(time.Time{}) {
					return data, nil
				}

				asString := data.(string)
				if asString == "" {
					return time.Time{}, nil
				}

				return time.Parse(time.RFC3339, asString)
			},
			func(
				f reflect.Type,
				t reflect.Type,
				data interface{}) (interface{}, error) {
				if f.Kind() != reflect.String {
					return data, nil
				}
				if t != reflect.TypeOf(log.InfoLevel) {
					return data, nil
				}

				asString := data.(string)
				if asString == "" {
					return log.InfoLevel, nil
				}

				return log.ParseLevel(asString)
			},
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)
	})
}
