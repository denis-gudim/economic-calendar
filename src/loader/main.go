package main

import (
	"economic-calendar/loader/investing/client"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	repository := &client.InvestingRepository{
		Source: &client.InvestingHttpClient{
			RetryCount: 10,
		},
		BatchSize:         4,
		DefaultLanguageId: 1,
	}

	//clt := &client.InvestingHttpClient{RetryCount: 10}

	//resp, err := clt.LoadEventDetailsHtml(1234, 1)
	//resp, err := clt.LoadCountriesHtml(1)
	//t := time.Date(2021, 9, 17, 0, 0, 0, 0, time.UTC)
	//items, err := repository.GetEventsSchedule(t, t)
	items, err := repository.GetCountries()

	if err == nil {
		fmt.Printf("OK nodes %v\n", len(items))
	} else {
		fmt.Printf("ERR %v\n", err)
	}
}
