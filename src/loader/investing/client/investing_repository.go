package client

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/investing/data"
	"economic-calendar/loader/investing/parsing"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

type InvestingDataEntry interface {
	GetId() int
	GetLanguageId() int
}

type InvestingHtmlSource interface {
	LoadEventsScheduleHtml(from, to time.Time, languageId int) (*goquery.Document, error)
	LoadEventDetailsHtml(eventId, languageId int) (*goquery.Document, error)
	LoadCountriesHtml(languageId int) (*goquery.Document, error)
}

type InvestingRepository struct {
	defaultLanguageId int
	batchSize         int
	source            InvestingHtmlSource
	logger            *log.Logger
}

func NewInvestingRepository(cnf app.Config, logger *log.Logger) *InvestingRepository {
	return &InvestingRepository{
		defaultLanguageId: cnf.Loading.DefaultLanguageId,
		batchSize:         cnf.Loading.BatchSize,
		source:            NewInvestingHttpClient(cnf),
		logger:            logger,
	}
}

func (repository *InvestingRepository) GetEventsSchedule(dateFrom, dateTo time.Time) (itemsMap map[int][]*data.InvestingScheduleRow, err error) {

	rows, err := repository.getItemsByLanguage(func(languageId int) ([]InvestingDataEntry, error) {
		return repository.getEventsScheduleByLanguage(languageId, dateFrom, dateTo)
	})

	if err != nil {
		return
	}

	count := len(data.InvestingLanguagesMap)
	itemsMap = make(map[int][]*data.InvestingScheduleRow, len(rows)/count)

	for _, row := range rows {

		row := row.(*data.InvestingScheduleRow)
		items, ok := itemsMap[row.Id]

		if !ok {
			items = make([]*data.InvestingScheduleRow, 0, count)
		}

		itemsMap[row.Id] = append(items, row)
	}

	return
}

func (repository *InvestingRepository) GetEventDetails(eventId int) (items []*data.InvestingCalendarEvent, err error) {

	rows, err := repository.getItemsByLanguage(func(languageId int) ([]InvestingDataEntry, error) {
		return repository.getEventDetailsByLanguage(languageId, eventId)
	})

	if err != nil {
		return
	}

	items = make([]*data.InvestingCalendarEvent, len(rows))

	for i, row := range rows {
		items[i] = row.(*data.InvestingCalendarEvent)
	}

	return
}

func (repository *InvestingRepository) GetCountries() (itemsMap map[int][]*data.InvestingCountry, err error) {

	rows, err := repository.getItemsByLanguage(func(languageId int) ([]InvestingDataEntry, error) {
		return repository.getCountriesByLanguage(languageId)
	})

	if err != nil {
		return
	}

	count := len(rows)/len(data.InvestingLanguagesMap) + 1
	itemsMap = make(map[int][]*data.InvestingCountry, len(data.InvestingLanguagesMap))

	for _, row := range rows {

		row := row.(*data.InvestingCountry)
		items, ok := itemsMap[row.Id]

		if !ok {
			items = make([]*data.InvestingCountry, 0, count)
		}

		itemsMap[row.Id] = append(items, row)
	}

	return
}

func (r *InvestingRepository) getEventsScheduleByLanguage(languageId int, dateFrom, dateTo time.Time) (items []InvestingDataEntry, err error) {

	html, err := r.source.LoadEventsScheduleHtml(dateFrom, dateTo, languageId)

	if err != nil {
		return
	}

	parser := parsing.NewInvestingScheduleParser()

	rows, err := parser.ParseScheduleHtml(html)

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

func (r *InvestingRepository) getEventDetailsByLanguage(languageId, eventId int) (event []InvestingDataEntry, err error) {
	html, err := r.source.LoadEventDetailsHtml(eventId, languageId)

	if err != nil {
		return
	}

	parser := parsing.NewInvestingCalendarEventParser()

	_event, err := parser.ParseCalendarEventHtml(html)

	if err != nil {
		return
	}

	_event.LanguageId = languageId

	return []InvestingDataEntry{_event}, nil
}

func (r *InvestingRepository) getCountriesByLanguage(languageId int) (items []InvestingDataEntry, err error) {
	html, err := r.source.LoadCountriesHtml(languageId)

	if err != nil {
		return
	}

	parser := &parsing.InvestingCountryParser{}

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

func (r *InvestingRepository) getItemsByLanguage(itemsGetter func(languageId int) ([]InvestingDataEntry, error)) (items []InvestingDataEntry, err error) {
	items, err = itemsGetter(r.defaultLanguageId)

	if err != nil {
		lang := data.InvestingLanguagesMap[r.defaultLanguageId]

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

	_itemsGetter := func(lang *data.InvestingLanguage) ([]InvestingDataEntry, error) {

		langItems, e := itemsGetter(lang.Id)

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

	count := len(data.InvestingLanguagesMap) - 1
	itemsChannel := make(chan []InvestingDataEntry, count)
	poolChannel := make(chan struct{}, batchSize)

	for languageId, language := range data.InvestingLanguagesMap {

		if languageId == r.defaultLanguageId {
			continue
		}

		poolChannel <- struct{}{}

		go func(lang *data.InvestingLanguage) {

			langItems, e := _itemsGetter(lang)

			if e != nil {
				r.logger.Errorf("items loading for language '%s' failed. %s", lang.Code, e.Error())
			}

			itemsChannel <- langItems
			<-poolChannel
		}(language)

	}

	for i := 0; i < count; i++ {
		items = append(items, <-itemsChannel...)
	}

	return
}
