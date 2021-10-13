package loading

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/data"
	"economic-calendar/loader/investing"
	"time"

	log "github.com/sirupsen/logrus"
)

type RefreshCalendarService struct {
	investingRepository     InvestingDataReciver
	eventScheduleRepository EventScheduleDataReciver
	logger                  *log.Logger
	config                  app.Config
	nextEventTime           time.Time
}

func NewRefreshCalendarService(cnf app.Config, logger *log.Logger) *RefreshCalendarService {
	return &RefreshCalendarService{
		investingRepository:     investing.NewInvestingRepository(cnf, logger),
		eventScheduleRepository: data.NewEventScheduleRepository(cnf),
		logger:                  logger,
		config:                  cnf,
	}
}

func (s *RefreshCalendarService) Refresh() {

	// nowTime := time.Now().UTC()

	// if nowTime.Before(s.nextEventTime) {

	// 	s.logger.WithFields(log.Fields{
	// 		"nowTime":       nowTime,
	// 		"nextEventTime": s.nextEventTime,
	// 	}).Debug("refresh skipped. next event time didn't come.")

	// 	return
	// }

	// fmtError := func(msg string, err error) error {
	// 	return xerrors.Errorf("events schedule refresh failed: %s: %w", msg, err)
	// }

	// s.logger.Info("events schedule refresh started...")

	// scheduleItems, err := s.eventScheduleRepository.GetByDates(nowTime, nowTime)

	// if err != nil {
	// 	err = fmtError("load countries", err)
	// 	s.logger.WithFields(log.Fields{ "nowTime": nowTime}).Error(err)
	// 	return
	// }

	// s.investingRepository
}
