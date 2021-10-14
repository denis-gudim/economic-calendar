package app

import (
	"os"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DB struct {
		ConnectionString string `yaml:"connectionString"`
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
