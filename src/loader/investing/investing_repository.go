package investing

import (
	"context"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/denis-gudim/economic-calendar/loader/app"
	log "github.com/sirupsen/logrus"
)

type InvestingDataEntry interface {
	GetId() int
	GetLanguageId() int
}

type InvestingHtmlSource interface {
	LoadEventsScheduleHtml(ctx context.Context, from, to time.Time, languageId int) (*goquery.Document, error)
	LoadEventDetailsHtml(ctx context.Context, eventId, languageId int) (*goquery.Document, error)
	LoadCountriesHtml(ctx context.Context, languageId int) (*goquery.Document, error)
}

type InvestingRepository struct {
	defaultLanguageId int
	batchSize         int
	source            InvestingHtmlSource
	logger            *log.Logger
}

func NewInvestingRepository(cnf *app.Config, logger *log.Logger, source InvestingHtmlSource) *InvestingRepository {
	return &InvestingRepository{
		defaultLanguageId: cnf.Loading.DefaultLanguageId,
		batchSize:         cnf.Loading.BatchSize,
		source:            source,
		logger:            logger,
	}
}

func (repository *InvestingRepository) GetEventsSchedule(ctx context.Context, dateFrom, dateTo time.Time) (itemsMap map[int][]*InvestingScheduleRow, err error) {

	rows, err := repository.getItemsByLanguage(ctx, func(ctx context.Context, languageId int) ([]InvestingDataEntry, error) {
		langItems, err := repository.GetEventsScheduleByLanguage(ctx, languageId, dateFrom, dateTo)

		if err != nil {
			return nil, err
		}

		items := make([]InvestingDataEntry, len(langItems))

		for i, item := range langItems {
			item.LanguageId = languageId
			items[i] = item
		}

		return items, nil
	})

	if err != nil {
		return
	}

	count := len(InvestingLanguagesMap)
	itemsMap = make(map[int][]*InvestingScheduleRow, len(rows)/count)

	for _, row := range rows {

		row := row.(*InvestingScheduleRow)
		items, ok := itemsMap[row.Id]

		if !ok {
			items = make([]*InvestingScheduleRow, 0, count)
		}

		itemsMap[row.Id] = append(items, row)
	}

	return
}

func (repository *InvestingRepository) GetEventDetails(ctx context.Context, eventId int) (items []*InvestingCalendarEvent, err error) {

	rows, err := repository.getItemsByLanguage(ctx, func(ctx context.Context, languageId int) ([]InvestingDataEntry, error) {
		return repository.getEventDetailsByLanguage(ctx, languageId, eventId)
	})

	if err != nil {
		return
	}

	items = make([]*InvestingCalendarEvent, len(rows))

	for i, row := range rows {
		items[i] = row.(*InvestingCalendarEvent)
	}

	return
}

func (repository *InvestingRepository) GetCountries(ctx context.Context) (itemsMap map[int][]*InvestingCountry, err error) {

	rows, err := repository.getItemsByLanguage(ctx, func(ctx context.Context, languageId int) ([]InvestingDataEntry, error) {
		return repository.getCountriesByLanguage(ctx, languageId)
	})

	if err != nil {
		return
	}

	count := len(rows)/len(InvestingLanguagesMap) + 1
	itemsMap = make(map[int][]*InvestingCountry, len(InvestingLanguagesMap))

	for _, row := range rows {

		row := row.(*InvestingCountry)
		items, ok := itemsMap[row.Id]

		if !ok {
			items = make([]*InvestingCountry, 0, count)
		}

		itemsMap[row.Id] = append(items, row)
	}

	return
}

func (r *InvestingRepository) GetEventsScheduleByLanguage(ctx context.Context, languageId int, dateFrom, dateTo time.Time) (items []*InvestingScheduleRow, err error) {

	html, err := r.source.LoadEventsScheduleHtml(ctx, dateFrom, dateTo, languageId)

	if err != nil {
		return
	}

	return NewInvestingScheduleParser().ParseScheduleHtml(html)
}

func (r *InvestingRepository) getEventDetailsByLanguage(ctx context.Context, languageId, eventId int) (event []InvestingDataEntry, err error) {
	html, err := r.source.LoadEventDetailsHtml(ctx, eventId, languageId)

	if err != nil {
		return
	}

	parser := NewInvestingCalendarEventParser()

	_event, err := parser.ParseCalendarEventHtml(html)

	if err != nil {
		return
	}

	_event.LanguageId = languageId

	return []InvestingDataEntry{_event}, nil
}

func (r *InvestingRepository) getCountriesByLanguage(ctx context.Context, languageId int) (items []InvestingDataEntry, err error) {
	html, err := r.source.LoadCountriesHtml(ctx, languageId)

	if err != nil {
		return
	}

	parser := &InvestingCountryParser{}

	rows, err := parser.ParseCountriesHtml(html)

	if err != nil {
		return
	}

	items = make([]InvestingDataEntry, len(rows))

	for i, row := range rows {
		row.LanguageId = languageId
		items[i] = row
	}

	return
}

func (r *InvestingRepository) getItemsByLanguage(ctx context.Context, itemsGetter func(ctx context.Context, languageId int) ([]InvestingDataEntry, error)) (items []InvestingDataEntry, err error) {
	items, err = itemsGetter(ctx, r.defaultLanguageId)

	if err != nil {
		lang := InvestingLanguagesMap[r.defaultLanguageId]

		log.Errorf("items loading for language '%s' failed. %s", lang.Code, err.Error())

		return
	}

	defLangItemsCount := len(items)

	if defLangItemsCount <= 0 {
		return
	}

	defaultLanguageItemsMap := make(map[int]InvestingDataEntry, defLangItemsCount)

	for _, item := range items {
		defaultLanguageItemsMap[item.GetId()] = item
	}

	_itemsGetter := func(lang *InvestingLanguage) ([]InvestingDataEntry, error) {

		langItems, e := itemsGetter(ctx, lang.Id)

		if e != nil {
			return nil, e
		}

		langItemsCount := len(langItems)

		if defLangItemsCount != langItemsCount {
			return nil, fmt.Errorf("items count not equals to default lang items %d/%d", langItemsCount, defLangItemsCount)
		}

		for _, item := range items {
			if _, ok := defaultLanguageItemsMap[item.GetId()]; !ok {
				return nil, fmt.Errorf("items have different keys with default items")
			}
		}

		return langItems, nil
	}

	batchSize := 1

	if r.batchSize > 0 {
		batchSize = r.batchSize
	}

	count := len(InvestingLanguagesMap) - 1
	itemsChannel := make(chan []InvestingDataEntry, count)
	poolChannel := make(chan struct{}, batchSize)

	for languageId, language := range InvestingLanguagesMap {

		if languageId == r.defaultLanguageId {
			continue
		}

		poolChannel <- struct{}{}

		go func(lang *InvestingLanguage) {

			langItems, e := _itemsGetter(lang)

			if e != nil {
				r.logger.Errorf("items loading for language '%s' failed. %s", lang.Code, e.Error())
			}

			itemsChannel <- langItems
			<-poolChannel
		}(language)

	}

	for i := 0; i < count; i++ {
		select {
		case batch := <-itemsChannel:
			items = append(items, batch...)
		case <-ctx.Done():
			return make([]InvestingDataEntry, 0), nil
		}
	}

	return
}
