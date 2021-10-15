package data

import (
	"context"
	"time"

	"github.com/denis-gudim/economic-calendar/api/app"
	"golang.org/x/xerrors"
)

type EventsRepository struct {
	baseRepository
}

func NewEventsRepository(cnf app.Config) *EventsRepository {
	r := EventsRepository{}
	r.config = cnf

	return &r
}

func (r *EventsRepository) GetScheduleByDates(ctx context.Context, from, to time.Time, langCode string) ([]Event, error) {
	const errFormat = "get schedule by dates failed: %s: %w"

	db, err := r.connectDB(ctx)

	if err != nil {
		return nil, xerrors.Errorf(errFormat, "connect to db", err)
	}

	defer db.Close()

	rows := make([]Event, 0, 128)

	err = db.SelectContext(ctx, &rows,
		`SELECT es.id, es.event_id, es.type, e.impact_level, c.code, es.timestamp_utc, est.title, es.actual, es.forecast, es.previous, e.unit
		 FROM event_schedule AS es JOIN events AS e 
		 ON e.id = es.event_id JOIN countries AS c
		 ON c.id = e.country_id JOIN event_schedule_translations AS est
		 ON es.id = est.event_schedule_id JOIN languages AS l
		 ON l.id = est.language_id and l.code = $1
		 WHERE es.timestamp_utc >= $2::timestamp AND es.timestamp_utc < $3::timestamp
		 ORDER BY es.timestamp_utc DESC`, langCode, from, to)

	if err != nil {
		return nil, xerrors.Errorf(errFormat, "execute select query", err)
	}

	return rows, nil
}

func (r *EventsRepository) GetEventById(ctx context.Context, eventId int, langCode string) ([]EventDetails, error) {
	const errFormat = "get event by id: %s: %w"

	db, err := r.connectDB(ctx)

	if err != nil {
		return nil, xerrors.Errorf(errFormat, "connect to db", err)
	}

	defer db.Close()

	rows := make([]EventDetails, 0, 1)

	err = db.SelectContext(ctx, &rows,
		`SELECT es.id, es.event_id, es.type, e.impact_level, c.code, es.timestamp_utc, et.title, es.actual, es.forecast, es.previous, et.overview, e.source, e.source_url, e.unit
		 FROM event_schedule AS es JOIN events AS e 
		 ON e.id = es.event_id AND e.id = $1 JOIN countries AS c
		 ON c.id = e.country_id JOIN event_translations AS et
		 ON e.id = et.event_id JOIN languages AS l
		 ON l.id = et.language_id and l.code = $2
		 ORDER BY es.timestamp_utc DESC
		 LIMIT 1`, eventId, langCode)

	if err != nil {
		return nil, xerrors.Errorf(errFormat, "execute select query", err)
	}

	return rows, nil
}

func (r *EventsRepository) GetHistoryById(ctx context.Context, eventId int) ([]EventRow, error) {
	const errFormat = "get history by id: %s: %w"

	db, err := r.connectDB(ctx)

	if err != nil {
		return nil, xerrors.Errorf(errFormat, "connect to db", err)
	}

	defer db.Close()

	rows := make([]EventRow, 0, 128)

	err = db.SelectContext(ctx, &rows,
		`SELECT id, event_id, timestamp_utc, actual, forecast, previous
		 FROM event_schedule
		 WHERE event_id = $1
		 ORDER BY timestamp_utc DESC`, eventId)

	if err != nil {
		return nil, xerrors.Errorf(errFormat, "execute select query", err)
	}

	return rows, nil

}
