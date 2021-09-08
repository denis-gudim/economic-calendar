package main

import (
	"economic-calendar/loader/client"
	"fmt"
	"time"
)

func main() {
	clt := &client.InvestingHttpClient{RetryCount: 10}
	resp, err := clt.LoadEventsScheduleHtml(time.Now(), time.Now(), 1)
	//resp, err := clt.LoadEventDetailsHtml(1234, 1)
	//resp, err := clt.LoadCountriesHtml(1)

	if err == nil {
		fmt.Println(resp.Find("table").Html())
		//fmt.Println(resp.Find("ul.countryOption").Html())
		//fmt.Println(resp.Find("#leftColumn").Html())
	} else {
		fmt.Println(err)
	}
}
