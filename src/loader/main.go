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

	if err := cnf.Load(); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	ds := loading.NewDictionariesLoaderService(cnf, log.StandardLogger())

	if err := ds.Load(); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	hs := loading.NewHistoryLoaderService(cnf, log.StandardLogger())

	hs.Load()

}
