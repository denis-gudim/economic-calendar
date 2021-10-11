package loading

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/data"
	"economic-calendar/loader/investing"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
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

type EventsDataReciver interface {
	GetById(id int) (*data.Event, error)
	Save(e data.Event) error
}

type EventScheduleDataReciver interface {
	GetFirst() (*data.EventSchedule, error)
	GetByDates(from, to time.Time) ([]data.EventSchedule, error)
	Save(es data.EventSchedule) error
}

type LoaderService struct {
	investingRepository     InvestingDataReciver
	countriesRepository     CountriesDataReciver
	eventsRepository        EventsDataReciver
	eventScheduleRepository EventScheduleDataReciver
	logger                  *log.Logger
	config                  app.Config
}

func NewLoaderService(cnf app.Config, logger *log.Logger) *LoaderService {

	return &LoaderService{
		investingRepository:     investing.NewInvestingRepository(cnf, logger),
		countriesRepository:     data.NewCountriesRepository(cnf),
		eventsRepository:        data.NewEventsRepository(cnf),
		eventScheduleRepository: data.NewEventScheduleRepository(cnf),
		logger:                  logger,
		config:                  cnf,
	}
}

func (s *LoaderService) LoadCountries() error {

	fmtError := func(err error) error {
		return xerrors.Errorf("countries dictionary loading error: %w", err)
	}

	s.logger.Info("countries dictionary loading started...")

	countries, err := s.countriesRepository.GetAll()

	if err != nil {
		return fmtError(err)
	}

	invCountries, err := s.investingRepository.GetCountries()

	if err != nil {
		return fmtError(err)
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
			return fmtError(err)
		}
	}

	s.logger.Info("countries dictionary loading finished")

	return nil
}

func (s *LoaderService) LoadEventsHistory() {

	s.logger.Info("events history loading started...")

	s.logger.Info("events history loading finished")
}
