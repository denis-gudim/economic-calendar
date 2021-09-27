package data

import "time"

type EventSchedule struct {
	Id        int
	TimeStamp time.Time
	Actual    *float64
	Forecast  *float64
	Previous  *float64
	IsDone    bool
	Type      int
}
