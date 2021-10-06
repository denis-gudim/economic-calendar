package loading

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/data"
	"economic-calendar/loader/investing"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type InvestingDataReciver interface {
	GetEventsSchedule(dateFrom, dateTo time.Time) (map[int][]*investing.InvestingScheduleRow, error)
	GetEventDetails(eventId int) ([]*investing.InvestingCalendarEvent, error)
	GetCountries() (map[int][]*investing.InvestingCountry, error)
}

type CountriesDataReciver interface {
	GetAll() ([]data.Country, error)
	Save(c data.Country) error
}

type LoaderService struct {
	investingRepository InvestingDataReciver
	countriesRepository CountriesDataReciver
	logger              *log.Logger
	config              app.Config
}

func NewLoaderService(cnf app.Config, logger *log.Logger) *LoaderService {

	return &LoaderService{
		investingRepository: investing.NewInvestingRepository(cnf, logger),
		countriesRepository: data.NewCountriesRepository(cnf),
		logger:              logger,
		config:              cnf,
	}
}

func (s *LoaderService) LoadCountries() error {

	s.logger.Info("countries dictionary loading started...")

	countries, err := s.countriesRepository.GetAll()

	if err != nil {
		return fmt.Errorf("countries dictionary loading error: %w", err)
	}

	invCountries, err := s.investingRepository.GetCountries()

	if err != nil {
		return fmt.Errorf("countries dictionary loading error: %w", err)
	}

	for _, c := range countries {
		ic, ok := invCountries[c.Id]

		if !ok {
			continue
		}

		c.Translations = make(data.Translations, len(ic))

		for _, icl := range ic {
			c.Translations[icl.LanguageId] = icl.Title
		}

		err = s.countriesRepository.Save(c)

		if err != nil {
			return fmt.Errorf("countries dictionary loading error: %w", err)
		}
	}

	s.logger.Info("countries dictionary loading finished")

	return nil
}
