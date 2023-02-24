package data

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type EventsRepository struct {
	baseRepository
}

func NewEventsRepository(db *sql.DB) *EventsRepository {
	r := EventsRepository{}
	r.db = db
	return &r
}

func (r *EventsRepository) GetById(ctx context.Context, id int) (*Event, error) {
	fmtError := func(msg string, err error) error {
		return fmt.Errorf("get events by id ( id: %d): %s: %w", id, msg, err)
	}

	filter := func(b sq.SelectBuilder) sq.SelectBuilder {
		return b.Where(sq.Eq{"id": id})
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

func (r *EventsRepository) Save(ctx context.Context, e Event) error {
	fmtError := func(msg string, err error) error {
		return fmt.Errorf("save event failed: %s: %w", msg, err)
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmtError("create db transaction", err)
	}
	defer func() {
		if tx != nil && err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				err = rerr
			}
		}
	}()

	upsertQuery := r.initQueryBuilder().
		Insert("events").
		Columns("id", "country_id", "impact_level", "unit", "source", "source_url").
		Values(e.Id, e.CountryId, e.ImpactLevel, e.Unit, e.Source, e.SourceUrl).
		Suffix("ON CONFLICT (id) DO").
		SuffixExpr(
			sq.Update(" ").
				Set("country_id", e.CountryId).
				Set("impact_level", e.ImpactLevel).
				Set("unit", e.Unit).
				Set("source", e.Source).
				Set("source_url", e.SourceUrl))

	_, err = upsertQuery.RunWith(tx).ExecContext(ctx)

	if err != nil {
		return fmtError("execute upsert event query", err)
	}

	deleteQuery := r.initQueryBuilder().
		Delete("event_translations").
		Where(sq.Eq{"event_id": e.Id})

	_, err = deleteQuery.RunWith(tx).ExecContext(ctx)

	if err != nil {
		return fmtError("execute delete event translations query", err)
	}

	for langId, title := range e.TitleTranslations {

		overview := e.OverviewTranslations[langId]

		insertQuery := r.initQueryBuilder().
			Insert("event_translations").
			Columns("event_id", "language_id", "title", "overview").
			Values(e.Id, langId, title, overview)

		_, err = insertQuery.RunWith(tx).ExecContext(ctx)

		if err != nil {
			return fmtError("execute insert event translation query", err)
		}
	}

	err = tx.Commit()

	if err != nil {
		return fmtError("commit transaction", err)
	}

	return nil
}

func (r *EventsRepository) getWithFilter(ctx context.Context, filter func(b sq.SelectBuilder) sq.SelectBuilder, fmtError func(suf string, err error) error) (events []Event, err error) {

	events = make([]Event, 0, 16)

	query := r.initQueryBuilder().
		Select("e.*, et.language_id, et.title, et.overview").
		From("events e").
		LeftJoin("event_translations et ON e.id = et.event_id").
		OrderBy("e.id")

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
		title     *string
		overview  *string
		curr      Event
		prevId    int
		titles    Translations
		overviews Translations
	)

	for rows.Next() {
		err = rows.Scan(
			&curr.Id,
			&curr.CountryId,
			&curr.ImpactLevel,
			&curr.Unit,
			&curr.Source,
			&curr.SourceUrl,
			&langId,
			&title,
			&overview,
		)

		if err != nil {
			return nil, fmtError("scan row", err)
		}

		if curr.Id != prevId {
			titles = Translations{}
			overviews = Translations{}
			curr.TitleTranslations = titles
			curr.OverviewTranslations = overviews
			events = append(events, curr)
			prevId = curr.Id
		}

		if title != nil {
			titles[*langId] = *title
		}

		if overview != nil {
			overviews[*langId] = *overview
		}
	}

	return
}
