package loading

import (
	"github.com/denis-gudim/economic-calendar/loader"
	log "github.com/sirupsen/logrus"
)

type RefreshCalendarService struct {
	investingRepository     InvestingDataReciver
	eventScheduleRepository EventScheduleDataReciver
	logger                  *log.Logger
	config                  *loader.Config
}

func NewRefreshCalendarService(cnf *loader.Config,
	logger *log.Logger,
	investingRepository InvestingDataReciver,
	eventScheduleRepository EventScheduleDataReciver) *RefreshCalendarService {

	return &RefreshCalendarService{
		investingRepository:     investingRepository,
		eventScheduleRepository: eventScheduleRepository,
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
	// 	return fmt.Errorf("events schedule refresh failed: %s: %w", msg, err)
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
