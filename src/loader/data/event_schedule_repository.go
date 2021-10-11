package data

import (
	"economic-calendar/loader/app"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/xerrors"
)

type EventScheduleRepository struct {
	baseRepository
}

func NewEventScheduleRepository(cnf app.Config) *EventScheduleRepository {
	r := EventScheduleRepository{}
	r.ConnectionString = cnf.DB.ConnectionString
	return &r
}

func (r *EventScheduleRepository) GetFirst() (es *EventSchedule, err error) {
	fmtError := func(text string, err error) error {
		return xerrors.Errorf("get first events schedule: %s: %w", text, err)
	}

	filter := func(b sq.SelectBuilder) {
		b.Limit(1)
	}

	res, err := r.getWithFilter(filter, fmtError)

	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

func (r *EventScheduleRepository) GetByDates(from, to time.Time) (events []EventSchedule, err error) {

	fromVal := from.Format("2006-01-02")
	toVal := to.Format("2006-01-02")

	fmtError := func(text string, err error) error {
		return xerrors.Errorf("get by dates ( from: %s, to: %s ) events schedule: %s: %w", fromVal, toVal, text, err)
	}

	filter := func(b sq.SelectBuilder) {
		b.Where(sq.And{
			sq.GtOrEq{"timestamp_utc": fromVal},
			sq.Lt{"timestamp_utc": toVal},
		})
	}

	return r.getWithFilter(filter, fmtError)
}

func (r *EventScheduleRepository) Save(es EventSchedule) error {

	fmtError := func(msg string, err error) error {
		return xerrors.Errorf("save event schedule failed: %s: %w", msg, err)
	}

	db, err := r.createConnection()

	if err != nil {
		return fmtError("create db connection", err)
	}

	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		return fmtError("create db transaction", err)
	}

	defer func() {
		if tx != nil && err != nil {
			tx.Rollback()
		}
	}()

	upsertQuery := r.initQueryBuilder().
		Insert("event_schedule").
		Columns("id", "timestamp_utc", "actual", "forecast", "previous", "done", "type", "event_id").
		Values(es.Id, es.TimeStamp, es.Actual, es.Forecast, es.Previous, es.IsDone, es.Type, es.EventId).
		Suffix("ON CONFLICT (id) DO").
		SuffixExpr(
			sq.Update(" ").
				Set("timestamp_utc", es.TimeStamp).
				Set("actual", es.Actual).
				Set("forecast", es.Forecast).
				Set("previous", es.Previous).
				Set("done", es.IsDone).
				Set("type", es.Type).
				Set("event_id", es.EventId))

	_, err = upsertQuery.RunWith(tx).Exec()

	if err != nil {
		return fmtError("execute upsert query", err)
	}

	deleteQuery := r.initQueryBuilder().
		Delete("event_schedule_translations").
		Where(sq.Eq{"event_schedule_id": es.Id})

	_, err = deleteQuery.RunWith(tx).Exec()

	if err != nil {
		return fmtError("execute delete translations query", err)
	}

	for langId, title := range es.Translations {

		insertQuery := r.initQueryBuilder().
			Insert("event_schedule_translations").
			Columns("event_schedule_id", "language_id", "title").
			Values(es.Id, langId, title)

		_, err = insertQuery.RunWith(tx).Exec()

		if err != nil {
			return fmtError("execute insert translation query", err)
		}
	}

	err = tx.Commit()

	if err != nil {
		return fmtError("commit transaction", err)
	}

	return nil
}

func (r *EventScheduleRepository) getWithFilter(filter func(b sq.SelectBuilder), fmtError func(suf string, err error) error) (events []EventSchedule, err error) {

	db, err := r.createConnection()

	if err != nil {
		return nil, fmtError("create db connection", err)
	}

	defer db.Close()

	events = make([]EventSchedule, 0, 256)

	query := r.initQueryBuilder().
		Select("es.*, et.language_id, et.title").
		From("event_schedule es").
		LeftJoin("event_schedule_translations est ON es.id = est.event_schedule_id").
		OrderBy("es.id")

	if filter != nil {
		filter(query)
	}

	rows, err := query.RunWith(db).Query()

	if err != nil {
		return nil, fmtError("execute select query", err)
	}

	var (
		langId    *int
		langTitle *string
		curr      EventSchedule
		prevId    int
		trans     Translations
	)

	for rows.Next() {
		err = rows.Scan(
			&curr.Id,
			&curr.TimeStamp,
			&curr.Actual,
			&curr.Forecast,
			&curr.Previous,
			&curr.IsDone,
			&curr.Type,
			&curr.EventId,
			&langId,
			&langTitle,
		)

		if err != nil {
			return nil, fmtError("scan row", err)
		}

		if curr.Id != prevId {
			trans = Translations{}
			curr.Translations = trans
			events = append(events, curr)
			prevId = curr.Id
		}

		if langId != nil {
			trans[*langId] = *langTitle
		}
	}

	return
}
