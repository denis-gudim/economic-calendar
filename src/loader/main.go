package main

import (
	"economic-calendar/loader/client"
	"fmt"
	"time"
)

func main() {
	clt := &client.InvestingHttpClient{RetryCount: 10}

	//resp, err := clt.LoadEventDetailsHtml(1234, 1)
	//resp, err := clt.LoadCountriesHtml(1)
	for languageId, language := range client.InvestingLanguagesMap {

		resp, err := clt.LoadEventsScheduleHtml(time.Now(), time.Now(), languageId)

		if err == nil {
			fmt.Printf("%s OK nodes %v\n", language.Code, len(resp.Nodes))
		} else {
			fmt.Printf("%s ERR %v\n", language.Code, err)
		}
	}
}
