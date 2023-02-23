package loading

import (
	"time"

	"github.com/denis-gudim/economic-calendar/loader"
	"github.com/denis-gudim/economic-calendar/loader/data"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/xerrors"
)

const daySec = 24 * 60 * 60

type HistoryLoaderService struct {
	investingRepository     InvestingDataReciver
	countriesRepository     CountriesDataReciver
	eventsRepository        EventsDataReciver
	eventScheduleRepository EventScheduleDataReciver
	logger                  *log.Logger
	config                  *loader.Config
	countriesMap            map[string]int
}

func NewHistoryLoaderService(cnf *loader.Config,
	logger *log.Logger,
	investingRepository InvestingDataReciver,
	countriesRepository CountriesDataReciver,
	eventsRepository EventsDataReciver,
	eventScheduleRepository EventScheduleDataReciver) *HistoryLoaderService {

	return &HistoryLoaderService{
		investingRepository:     investingRepository,
		countriesRepository:     countriesRepository,
		eventsRepository:        eventsRepository,
		eventScheduleRepository: eventScheduleRepository,
		logger:                  logger,
		config:                  cnf,
	}
}

func (s *HistoryLoaderService) Load(ctx context.Context) {

	fmtError := func(msg string, err error) error {
		return xerrors.Errorf("events schedule loading failed: %s: %w", msg, err)
	}

	s.logger.Info("events history loading started...")

	err := s.fillCountriesMap(ctx)

	if err != nil {
		s.logger.Error(fmtError("fill countries map", err))
	}

	from, to, err := s.getHistoryLoadingDates(ctx)

	if err != nil {
		s.logger.Error(fmtError("loading dates calculation", err))
		return
	}

	_ctx, cancelFunc := context.WithCancel(ctx)

	defer cancelFunc()

	out1, errc1 := s.loadInvestingSchedule(_ctx, from, to)
	out2, errc2 := s.loadInvestingEvents(_ctx, out1)

	for item := range out2 {
		err = s.eventScheduleRepository.Save(ctx, item)

		if err != nil {
			s.logger.Error(fmtError("save loaded schedule items", err))
			break
		}

		s.logger.Infof("new event schedule row stored to database: id = %d, eventId = %d", item.Id, item.EventId)
	}

	select {
	case err := <-errc1:
		if err != nil {
			s.logger.Error(fmtError("load investing schedule", err))
		}
	case err := <-errc2:
		if err != nil {
			s.logger.Error(fmtError("load investing events", err))
		}
	case <-_ctx.Done():
		{
			s.logger.Info("events history loading canceled")
		}
	default:
		s.logger.Info("events history loading finished")
	}
}

func (s *HistoryLoaderService) fillCountriesMap(ctx context.Context) error {
	countries, err := s.countriesRepository.GetAll(ctx)

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

func (s *HistoryLoaderService) getHistoryLoadingDates(ctx context.Context) (from, to time.Time, err error) {

	fromRow, err := s.eventScheduleRepository.GetFirst(ctx, true)

	if err != nil {
		return
	}

	toRow, err := s.eventScheduleRepository.GetFirst(ctx, false)

	if err != nil {
		return
	}

	if fromRow != nil {
		from = time.Unix(fromRow.TimeStamp.Unix()/daySec*daySec, 0)
	}

	if toRow != nil {
		to = time.Unix(toRow.TimeStamp.Unix()/daySec*daySec, 0)
	}

	return from, to, nil
}

func (s *HistoryLoaderService) loadInvestingSchedule(ctx context.Context, from, to time.Time) (<-chan data.EventSchedule, <-chan error) {

	out := make(chan data.EventSchedule, 1024)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		date := time.Unix((time.Now().UTC().Unix()/daySec+int64(s.config.Loading.ToDays))*daySec, 0)

		for !date.Before(s.config.Loading.FromTime) {

			if !date.After(from) || !date.Before(to) {

				batch, err := s.investingRepository.GetEventsSchedule(ctx, date, date)

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
						IsDone:            langItem.IsDone(time.Now().UTC()),
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

			event, err := s.eventsRepository.GetById(ctx, scheduleRow.EventId)

			if err != nil {
				errc <- err
				return
			}

			if event == nil {

				translations, err := s.investingRepository.GetEventDetails(ctx, scheduleRow.EventId)

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

				err = s.eventsRepository.Save(ctx, newEvent)

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
