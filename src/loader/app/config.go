package app

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		ConnectionString string `yaml:"connectionString"`
	}
	Loading struct {
		DefaultLanguageId int       `yaml:"defaultLanguageId"`
		RetryCount        int       `yaml:"retryCount"`
		BatchSize         int       `yaml:"batchSize"`
		FromTime          time.Time `yaml:"fromTime"`
		ToDays            int       `yaml:"toDays"`
	}
	Logging struct {
		Level log.Level `yaml:"level"`
	}
	Scheduler struct {
		HistoryExpression string `yaml:"historyExpression"`
		RefreshExpression string `yaml:"refreshExpression"`
	}
}

func (cnf *Config) Load() error {
	err := cnf.loadYamlConfig()

	if err != nil {
		return err
	}

	return cnf.loadEnvConfig()
}

func (cnf *Config) loadYamlConfig() error {
	file, err := os.Open("config.yaml")

	if err != nil {
		return xerrors.Errorf("load yaml config failed: %w", err)
	}

	defer file.Close()

	err = yaml.NewDecoder(file).Decode(cnf)

	if err != nil {
		return xerrors.Errorf("load yaml config failed: %w", err)
	}

	return nil
}

func (cnf *Config) loadEnvConfig() error {
	return nil
}
