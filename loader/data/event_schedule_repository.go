package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type EventScheduleRepository struct {
	baseRepository
}

func NewEventScheduleRepository(db *sql.DB) *EventScheduleRepository {
	r := EventScheduleRepository{}
	r.db = db
	return &r
}

func (r *EventScheduleRepository) GetFirst(ctx context.Context, done bool) (es *EventSchedule, err error) {
	fmtError := func(text string, err error) error {
		return fmt.Errorf("get first events schedule: %s: %w", text, err)
	}

	filter := func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{"done": done}).Limit(1)
	}

	res, err := r.getWithFilter(ctx, filter, fmtError)

	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return &res[0], nil
}

func (r *EventScheduleRepository) GetByDates(ctx context.Context, from, to time.Time) (events []EventSchedule, err error) {
	fromVal := from.Format("2006-01-02")
	toVal := to.Format("2006-01-02")

	fmtError := func(text string, err error) error {
		return fmt.Errorf("get by dates ( from: %s, to: %s ) events schedule: %s: %w", fromVal, toVal, text, err)
	}

	filter := func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.And{
			sq.GtOrEq{"timestamp_utc": fromVal},
			sq.Lt{"timestamp_utc": toVal},
		})
	}

	return r.getWithFilter(ctx, filter, fmtError)
}

func (r *EventScheduleRepository) Save(ctx context.Context, es EventSchedule) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("create db transaction error: %w", err)
	}
	defer func() {
		if tx != nil && err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				err = rerr
			}
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
	_, err = upsertQuery.RunWith(tx).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("execute upsert query error: %w", err)
	}

	deleteQuery := r.initQueryBuilder().
		Delete("event_schedule_translations").
		Where(sq.Eq{"event_schedule_id": es.Id})
	_, err = deleteQuery.RunWith(tx).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("execute delete translations query error: %w", err)
	}

	for langId, title := range es.TitleTranslations {
		insertQuery := r.initQueryBuilder().
			Insert("event_schedule_translations").
			Columns("event_schedule_id", "language_id", "title").
			Values(es.Id, langId, title)
		_, err = insertQuery.RunWith(tx).ExecContext(ctx)
		if err != nil {
			return fmt.Errorf("execute insert translation query: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit transaction error: %w", err)
	}

	return nil
}

func (r *EventScheduleRepository) getWithFilter(ctx context.Context, filter func(b sq.SelectBuilder) sq.SelectBuilder, fmtError func(suf string, err error) error) (events []EventSchedule, err error) {

	events = make([]EventSchedule, 0, 256)

	query := r.initQueryBuilder().
		Select("es.*, est.language_id, est.title").
		From("event_schedule es").
		LeftJoin("event_schedule_translations est ON es.id = est.event_schedule_id").
		OrderBy("es.id")

	if filter != nil {
		query = filter(query)
	}

	rows, err := query.RunWith(r.db).QueryContext(ctx)
	if err != nil {
		return nil, fmtError("execute select query", err)
	}
	defer rows.Close()

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
			curr.TitleTranslations = trans
			events = append(events, curr)
			prevId = curr.Id
		}

		if langId != nil {
			trans[*langId] = *langTitle
		}
	}

	return
}
