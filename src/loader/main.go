package main

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/data"
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

	cnf := app.Config{}
	cnf.Load()

	//	dbTest(&cnf)

	invTest(&cnf)

}

func dbTest(cnf *app.Config) {
	lr := data.LanguagesRepository{}
	lr.ConnectionString = cnf.DB.ConnectionString
	languages, e := lr.GetAll()

	cr := data.CountriesRepository{}
	cr.ConnectionString = cnf.DB.ConnectionString
	countries, e := cr.GetAll()

	for _, c := range countries {
		c.Name = c.Name + "_"

		for _, v := range languages {
			c.Translations[v.Code] = c.Name + "_" + v.NativeName
		}

		e = cr.Save(&c)

		if e != nil {
			fmt.Println(e)
		}
	}

	countries, e = cr.GetAll()

	if e == nil {
		fmt.Printf("countries : %v", len(countries))
	} else {
		fmt.Print(e)
	}
}

func invTest(cnf *app.Config) {
	repository := &client.InvestingRepository{
		Source: &client.InvestingHttpClient{
			RetryCount: 10,
		},
		BatchSize:         4,
		DefaultLanguageId: 1,
		Logger:            log.New(),
	}

	//t := time.Date(2021, 9, 17, 0, 0, 0, 0, time.UTC)
	//items, err := repository.GetEventsSchedule(t, t)

	//items, err := repository.GetEventDetails(1042)

	items, err := repository.GetCountries()

	for _, c := range items[1] {
		fmt.Printf("%v %v\n", c.Id, c.Title)
	}

	if err == nil {
		fmt.Printf("OK nodes %v\n", len(items))
	} else {
		fmt.Printf("ERR %v\n", err)
	}
}
