package main

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/loading"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	cnf := app.Config{}
	cnf.Load()

	l := loading.NewLoaderService(cnf, log.StandardLogger())

	if err := l.LoadCountries(); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

}
