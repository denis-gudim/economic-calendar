package data

import "time"

type EventRow struct {
	Id        int       `json:"id"`
	EventId   int       `db:"event_id" json:"eventId"`
	Timestamp time.Time `db:"timestamp_utc" json:"timestamp"`
	Actual    *float64  `json:"actual"`
	Forecast  *float64  `json:"forecast"`
	Previous  *float64  `json:"previous"`
}
