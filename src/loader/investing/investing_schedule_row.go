package investing

import (
	"time"
)

type ScheduleEventType int

const (
	Index ScheduleEventType = iota
	Speech
	PreliminaryRelease
	Report
	RetrievingData
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
	Type         ScheduleEventType
}

func (row *InvestingScheduleRow) GetId() int {
	return row.Id
}

func (row *InvestingScheduleRow) GetLanguageId() int {
	return row.LanguageId
}
