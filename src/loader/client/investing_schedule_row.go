package client

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
