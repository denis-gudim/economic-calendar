package app

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		ConnectionString string `yaml:"connectionString"`
	}
	Loading struct {
		RetryCount int `yaml:"retryCount"`
		BatchSize  int `yaml:"batchSize"`
	}
}

func (cnf *Config) Load() {
	cnf.loadYamlConfig()
	cnf.loadEnvConfig()
}

func (cnf *Config) loadYamlConfig() {
	file, err := os.Open("config.yaml")

	if err != nil {
		cnf.processError(err)
	}

	defer file.Close()

	err = yaml.NewDecoder(file).Decode(cnf)

	if err != nil {
		cnf.processError(err)
	}
}

func (cnf *Config) loadEnvConfig() {

}

func (cnf *Config) processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
