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

func (r *InvestingScheduleRow) GetId() int {
	return r.Id
}

func (r *InvestingScheduleRow) GetLanguageId() int {
	return r.LanguageId
}

func (r *InvestingScheduleRow) IsDone(time time.Time) bool {
	if r.Type == Index {
		return r.Actual != nil
	}

	return time.After(r.TimeStamp)
}
