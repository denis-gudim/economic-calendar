package data

import "time"

type ScheduleRow struct {
	Id          int       `json:"id"`
	EventId     int       `db:"event_id" json:"eventId"`
	Type        int       `json:"type"`
	ImpactLevel int       `db:"impact_level" json:"impactLevel"`
	Code        string    `json:"code"`
	Timestamp   time.Time `db:"timestamp_utc" json:"timestamp"`
	Unit        string    `json:"unit"`
	Title       string    `json:"title"`
	Actual      *float64  `json:"actual"`
	Forecast    *float64  `json:"forecast"`
	Previous    *float64  `json:"previous"`
}

type ScheduleEventDetails struct {
	ScheduleRow
	Overview  string `json:"overview"`
	Source    string `json:"source"`
	SourceUrl string `db:"source_url" json:"sourceUrl"`
}
