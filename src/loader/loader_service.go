package main

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/data"
	investing_data "economic-calendar/loader/investing/data"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type InvestingDataReciver interface {
	GetEventsSchedule(dateFrom, dateTo time.Time) ([]*investing_data.InvestingScheduleRow, error)
	GetEventDetails(eventId int) ([]*investing_data.InvestingCalendarEvent, error)
	GetCountries() (map[int][]*investing_data.InvestingCountry, error)
}

type CountriesDataReciver interface {
	GetAll() ([]data.Country, error)
}

type LoaderService struct {
	investingRepository InvestingDataReciver
	countriesRepository CountriesDataReciver
	Logger              *log.Logger
	Config              app.Config
}

func NewLoaderService() {

}

func (s *LoaderService) LoadCountries() error {

	s.Logger.Info("load countries started...")

	invCountries, err := s.investingRepository.GetCountries()

	if err != nil {
		return fmt.Errorf("load countries failed: %w", err)
	}

	countries, err := s.countriesRepository.GetAll()

	if err != nil {
		return fmt.Errorf("load countries failed: %w", err)
	}

	for _, c := range countries {
		ic, ok := invCountries[*c.InvestingId]

		if !ok {
			continue
		}

		c.Translations = make(data.Translations, len(ic))

		for _, icl := range ic {
			lang := investing_data.InvestingLanguagesMap[icl.LanguageId]

			c.Translations[lang.Code]
		}
	}

	return nil
}
