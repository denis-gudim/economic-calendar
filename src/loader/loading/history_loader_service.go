package loading

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/data"
	"economic-calendar/loader/investing"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/xerrors"
)

type HistoryLoaderService struct {
	investingRepository     InvestingDataReciver
	countriesRepository     CountriesDataReciver
	eventsRepository        EventsDataReciver
	eventScheduleRepository EventScheduleDataReciver
	logger                  *log.Logger
	config                  app.Config
	countriesMap            map[string]int
}

func NewHistoryLoaderService(cnf app.Config, logger *log.Logger) *HistoryLoaderService {

	return &HistoryLoaderService{
		investingRepository:     investing.NewInvestingRepository(cnf, logger),
		countriesRepository:     data.NewCountriesRepository(cnf),
		eventsRepository:        data.NewEventsRepository(cnf),
		eventScheduleRepository: data.NewEventScheduleRepository(cnf),
		logger:                  logger,
		config:                  cnf,
	}
}

func (s *HistoryLoaderService) Load() {

	fmtError := func(msg string, err error) error {
		return xerrors.Errorf("events schedule loading failed: %s: %w", msg, err)
	}

	s.logger.Info("events history loading started...")

	err := s.fillCountriesMap()

	if err != nil {
		s.logger.Error(fmtError("fill countries map", err))
	}

	startDate, err := s.getHistoryLoadingStartDate()

	if err != nil {
		s.logger.Error(fmtError("start loading date calculation", err))
		return
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	defer cancelFunc()

	out1, errc1 := s.loadInvestingSchedule(ctx, startDate)
	out2, errc2 := s.loadInvestingEvents(ctx, out1)

	for item := range out2 {
		err = s.eventScheduleRepository.Save(item)

		if err != nil {
			s.logger.Error(fmtError("save loaded schedule items", err))
			break
		}

		s.logger.Infof("new event schedule row stored to database: id = %d, eventId = %d", item.Id, item.EventId)
	}

	select {
	case err := <-errc1:
		s.logger.Error(fmtError("load investing schedule", err))
	case err := <-errc2:
		s.logger.Error(fmtError("load investing events", err))
	default:
		s.logger.Info("events history loading finished")
	}
}

func (s *HistoryLoaderService) fillCountriesMap() error {
	countries, err := s.countriesRepository.GetAll()

	if err != nil {
		return err
	}

	cm := make(map[string]int, len(countries))

	for _, c := range countries {
		cm[c.Name] = c.Id
	}

	s.countriesMap = cm

	return nil
}

func (s *HistoryLoaderService) getHistoryLoadingStartDate() (t time.Time, err error) {
	frow, err := s.eventScheduleRepository.GetFirst()

	if err != nil {
		return
	}

	if frow == nil {
		return time.Now().AddDate(0, 0, s.config.Loading.ToDays), nil
	}

	return frow.TimeStamp, nil
}

func (s *HistoryLoaderService) loadInvestingSchedule(ctx context.Context, date time.Time) (<-chan data.EventSchedule, <-chan error) {

	out := make(chan data.EventSchedule, 1024)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for !date.Before(s.config.Loading.FromTime) {

			batch, err := s.investingRepository.GetEventsSchedule(date, date)

			if err != nil {
				errc <- err
				break
			}

			s.logger.Infof("events schedule history batch loaded: date = %s, count = %d", date, len(batch))

			for rowId, translations := range batch {

				langItem := translations[0]

				if len(translations) == 0 {
					errc <- xerrors.Errorf("translations list is empty")
					return
				}

				newScheduleRow := data.EventSchedule{
					Id:                rowId,
					TimeStamp:         langItem.TimeStamp,
					Actual:            langItem.Actual,
					Forecast:          langItem.Forecast,
					Previous:          langItem.Previous,
					IsDone:            langItem.IsDone(time.Now()),
					Type:              int(langItem.Type),
					EventId:           langItem.EventId,
					TitleTranslations: data.Translations{},
				}

				for _, langItem = range translations {
					newScheduleRow.TitleTranslations[langItem.LanguageId] = langItem.Title
				}

				select {
				case out <- newScheduleRow:
				case <-ctx.Done():
					return
				}
			}

			date = date.AddDate(0, 0, -1)
		}
	}()

	return out, errc
}

func (s *HistoryLoaderService) loadInvestingEvents(ctx context.Context, in <-chan data.EventSchedule) (<-chan data.EventSchedule, <-chan error) {

	out := make(chan data.EventSchedule, 1024)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for scheduleRow := range in {

			event, err := s.eventsRepository.GetById(scheduleRow.EventId)

			if err != nil {
				errc <- err
				return
			}

			if event == nil {

				translations, err := s.investingRepository.GetEventDetails(scheduleRow.EventId)

				s.logger.Infof("event details loaded from source: eventId = %d", scheduleRow.EventId)

				if err != nil {
					errc <- err
					return
				}

				if len(translations) == 0 {
					errc <- xerrors.Errorf("translations list is empty")
					return
				}

				langItem := translations[0]

				countryId, ok := s.countriesMap[langItem.Country]

				if !ok {
					errc <- xerrors.Errorf("country with name '%s' not found in map", langItem.Country)
					return
				}

				newEvent := data.Event{
					Id:                   scheduleRow.EventId,
					CountryId:            countryId,
					ImpactLevel:          langItem.Sentiment,
					Unit:                 langItem.Unit,
					Source:               langItem.Source,
					SourceUrl:            langItem.SourceUrl,
					TitleTranslations:    data.Translations{},
					OverviewTranslations: data.Translations{},
				}

				for _, langItem = range translations {
					newEvent.TitleTranslations[langItem.LanguageId] = langItem.Title
					newEvent.OverviewTranslations[langItem.LanguageId] = langItem.Overview
				}

				err = s.eventsRepository.Save(newEvent)

				if err != nil {
					errc <- err
					return
				}

				s.logger.Infof("new event details stored to database: id = %d", newEvent.Id)
			}

			select {
			case out <- scheduleRow:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, errc
}
