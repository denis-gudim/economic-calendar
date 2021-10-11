package data

import (
	"economic-calendar/loader/app"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/xerrors"
)

type EventsRepository struct {
	baseRepository
}

func NewEventsRepository(cnf app.Config) *EventsRepository {
	r := EventsRepository{}
	r.ConnectionString = cnf.DB.ConnectionString
	return &r
}

func (r *EventsRepository) GetById(id int) (*Event, error) {

	fmtError := func(msg string, err error) error {
		return xerrors.Errorf("get events by id ( id: %d): %s: %w", id, msg, err)
	}

	filter := func(b sq.SelectBuilder) {
		b.Where(sq.Eq{"id": id})
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

func (r *EventsRepository) Save(e Event) error {

	fmtError := func(msg string, err error) error {
		return xerrors.Errorf("save event failed: %s: %w", msg, err)
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

	_, err = upsertQuery.RunWith(tx).Exec()

	if err != nil {
		return fmtError("execute upsert event query", err)
	}

	deleteQuery := r.initQueryBuilder().
		Delete("event_translations").
		Where(sq.Eq{"event_id": e.Id})

	_, err = deleteQuery.RunWith(tx).Exec()

	if err != nil {
		return fmtError("execute delete event translations query", err)
	}

	for langId, title := range e.Translations {

		insertQuery := r.initQueryBuilder().
			Insert("event_translations").
			Columns("event_id", "language_id", "title").
			Values(e.Id, langId, title)

		_, err = insertQuery.RunWith(tx).Exec()

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

func (r *EventsRepository) getWithFilter(filter func(b sq.SelectBuilder), fmtError func(suf string, err error) error) (events []Event, err error) {

	db, err := r.createConnection()

	if err != nil {
		return nil, fmtError("create db connection", err)
	}

	defer db.Close()

	events = make([]Event, 0, 16)

	query := r.initQueryBuilder().
		Select("e.*, et.language_id, et.title").
		From("events e").
		LeftJoin("event_translations et ON e.id = et.event_id").
		OrderBy("e.id")

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
		curr      Event
		prevId    int
		trans     Translations
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
