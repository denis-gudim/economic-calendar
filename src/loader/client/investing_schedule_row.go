package client

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ScheduleEventType int

const (
	Index              ScheduleEventType = iota
	Speech             ScheduleEventType = iota
	PreliminaryRelease ScheduleEventType = iota
	Report             ScheduleEventType = iota
	RetrievingData     ScheduleEventType = iota
)

type InvestingScheduleRow struct {
	EventId      int
	Id           int
	LanguageId   int
	CountryName  string
	Title        string
	TimeStamp    time.Time
	CurrencyCode string
	Sentiment    int
	Actual       *float64
	Forecast     *float64
	Previous     *float64
	Revised      *float64
	Type         ScheduleEventType
}

func parseScheduleHtml(html *goquery.Document) (rows []*InvestingScheduleRow, err error) {
	return nil, nil
}

func parseScheduleRowHtml(s *goquery.Selection) (row *InvestingScheduleRow, err error) {
	_, idAttrExists := s.Attr("id")
	_, eventIdAttrExists := s.Attr("event_attr_id")

	if idAttrExists && eventIdAttrExists {

		row = &InvestingScheduleRow{}

		row.EventId, err = parseAttrValueToInt(s, "event_attr_id")

		if err != nil {
			return
		}

	}

	return
}

var idRegEx *regexp.Regexp = regexp.MustCompile(`(\d+)`)

func parseScheduleRowId(s *goquery.Selection) (id int, err error) {

	idVal, err := getAttrValue(s, "id")

	if err != nil {
		return
	}

	idStr := idRegEx.FindString(idVal)

	if len(idStr) <= 0 {
		return 0, &ParsingError{
			Err: fmt.Errorf("id attribute has invalid value '%s'", idStr),
		}
	}

	id, err = strconv.Atoi(idStr)

	if err != nil {
		return 0, &ParsingError{Err: err}
	}

	return
}

func parseScheduleCountryName(s *goquery.Selection) (string, error) {

	flagCell := s.Find("span.flagCur")

	return getAttrValue(flagCell, "title")
}
