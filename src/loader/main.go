package main

import (
	"economic-calendar/loader/client"
	"fmt"
)

func main() {
	clt := &client.InvestingHttpClient{RetryCount: 10}
	//resp, err := clt.LoadEventsScheduleHtml(time.Now(), time.Now(), 1)
	//resp, err := clt.LoadEventDetailsHtml(1234, 1)
	resp, err := clt.LoadCountriesHtml(1)

	if err == nil {
		fmt.Println(resp)
	} else {
		fmt.Println(err)
	}
}
