package client

import (
	"economic-calendar/loader/investing/data"
	"economic-calendar/loader/investing/parsing"
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
	DefaultLanguageId int
	BatchSize         int
	Source            InvestingHtmlSource
}

func (repository *InvestingRepository) GetEventsSchedule(dateFrom, dateTo time.Time) (items []*data.InvestingScheduleRow, err error) {

	rows, err := repository.getItemsByLanguage(func(languageId int) ([]InvestingDataEntry, error) {
		return repository.getEventsScheduleByLanguage(languageId, dateFrom, dateTo)
	})

	if err != nil {
		return
	}

	items = make([]*data.InvestingScheduleRow, len(rows))

	for i, row := range rows {
		items[i] = row.(*data.InvestingScheduleRow)
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

func (repository *InvestingRepository) GetCountries() (items []*data.InvestingCountry, err error) {

	rows, err := repository.getItemsByLanguage(func(languageId int) ([]InvestingDataEntry, error) {
		return repository.getCountriesByLanguage(languageId)
	})

	if err != nil {
		return
	}

	items = make([]*data.InvestingCountry, len(rows))

	for i, row := range rows {
		items[i] = row.(*data.InvestingCountry)
	}

	return
}

func (repository *InvestingRepository) getEventsScheduleByLanguage(languageId int, dateFrom, dateTo time.Time) (items []InvestingDataEntry, err error) {

	html, err := repository.Source.LoadEventsScheduleHtml(dateFrom, dateTo, languageId)

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

func (repository *InvestingRepository) getEventDetailsByLanguage(languageId, eventId int) (event []InvestingDataEntry, err error) {
	html, err := repository.Source.LoadEventDetailsHtml(languageId, eventId)

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

func (repository *InvestingRepository) getCountriesByLanguage(languageId int) (items []InvestingDataEntry, err error) {
	html, err := repository.Source.LoadCountriesHtml(languageId)

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

func (repository *InvestingRepository) getItemsByLanguage(itemsGetter func(languageId int) ([]InvestingDataEntry, error)) (items []InvestingDataEntry, err error) {
	items, err = itemsGetter(repository.DefaultLanguageId)

	if err != nil {
		lang := data.InvestingLanguagesMap[repository.DefaultLanguageId]

		log.Errorf("items loading for language '%s' failed. %s", lang.Code, err.Error())

		return
	}

	defaultLanguageItemsCount := len(items)

	if defaultLanguageItemsCount <= 0 {
		return
	}

	defaultLanguageItemsMap := make(map[int]InvestingDataEntry, defaultLanguageItemsCount)

	for _, item := range items {
		defaultLanguageItemsMap[item.GetId()] = item
	}

	itemsComparer := func(items []InvestingDataEntry) bool {

		if defaultLanguageItemsCount != len(items) {
			return false
		}

		for _, item := range items {
			if _, ok := defaultLanguageItemsMap[item.GetId()]; !ok {
				return false
			}
		}

		return true
	}

	batchSize := 1

	if repository.BatchSize > 0 {
		batchSize = repository.BatchSize
	}

	count := len(data.InvestingLanguagesMap) - 1
	itemsChannel := make(chan []InvestingDataEntry, count)
	poolChannel := make(chan struct{}, batchSize)

	for languageId, language := range data.InvestingLanguagesMap {

		if languageId == repository.DefaultLanguageId {
			continue
		}

		poolChannel <- struct{}{}

		go func(lang data.InvestingLanguage) {

			languageItems, e := itemsGetter(lang.Id)

			if e != nil {
				log.Errorf("items loading for language '%s' failed. %s", lang.Code, e.Error())
			} else if !itemsComparer(languageItems) {
				log.Errorf("items loading for language '%s' failed. language items are different", lang.Code)
			}

			itemsChannel <- languageItems
			<-poolChannel
		}(language)

	}

	for i := 0; i < count; i++ {
		items = append(items, <-itemsChannel...)
	}

	return
}