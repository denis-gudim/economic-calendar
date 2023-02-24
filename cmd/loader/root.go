package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/denis-gudim/economic-calendar/loader"
	"github.com/denis-gudim/economic-calendar/loader/data"
	"github.com/denis-gudim/economic-calendar/loader/investing"
	"github.com/denis-gudim/economic-calendar/loader/loading"
	"github.com/go-co-op/gocron"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type CompositionRoot struct {
	logger    *logrus.Logger
	db        *sql.DB
	container *dig.Container
	cnf       *loader.Config
}

func NewCompositionRoot() (*CompositionRoot, error) {
	container := dig.New()

	cnf := &loader.Config{}
	if err := cnf.Load(); err != nil {
		return nil, fmt.Errorf("load config error: %w", err)
	}

	logger := logrus.StandardLogger()
	logger.SetLevel(cnf.Logging.Level)

	db, err := sql.Open("postgres", cnf.DB.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("connect to database error: %w", err)
	}

	err = container.Provide(func() *loader.Config {
		return cnf
	})
	if err != nil {
		return nil, err
	}
	err = container.Provide(func() *logrus.Logger {
		return logger
	})
	if err != nil {
		return nil, err
	}
	err = container.Provide(func() *sql.DB {
		return db
	})
	if err != nil {
		return nil, err
	}
	err = container.Provide(func(c *loader.Config) investing.InvestingHtmlSource {
		return investing.NewInvestingHttpClient(c)
	})
	if err != nil {
		return nil, err
	}
	err = container.Provide(func(c *loader.Config, logger *logrus.Logger, source investing.InvestingHtmlSource) loading.InvestingDataReciver {
		return investing.NewInvestingRepository(c, logger, source)
	})
	if err != nil {
		return nil, err
	}
	err = container.Provide(func(db *sql.DB) loading.CountriesDataReciver {
		return data.NewCountriesRepository(db)
	})
	if err != nil {
		return nil, err
	}
	err = container.Provide(func(db *sql.DB) loading.EventScheduleDataReciver {
		return data.NewEventScheduleRepository(db)
	})
	if err != nil {
		return nil, err
	}
	err = container.Provide(func(db *sql.DB) loading.EventsDataReciver {
		return data.NewEventsRepository(db)
	})
	if err != nil {
		return nil, err
	}
	err = container.Provide(func(db *sql.DB) loading.EventsDataReciver {
		return data.NewEventsRepository(db)
	})
	if err != nil {
		return nil, err
	}
	err = container.Provide(loading.NewDictionariesLoaderService)
	if err != nil {
		return nil, err
	}
	err = container.Provide(loading.NewHistoryLoaderService)
	if err != nil {
		return nil, err
	}
	err = container.Provide(loading.NewRefreshCalendarService)
	if err != nil {
		return nil, err
	}
	err = container.Provide(NewHealtz)
	if err != nil {
		return nil, err
	}

	return &CompositionRoot{
		db:        db,
		logger:    logger,
		container: container,
		cnf:       cnf,
	}, nil
}

func (r *CompositionRoot) GetContainer() *dig.Container {
	return r.container
}

func (r *CompositionRoot) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

func (r *CompositionRoot) InitSchedule(ctx context.Context, s *gocron.Scheduler) error {

	err := r.container.Invoke(func(cnf *loader.Config, srv *loading.HistoryLoaderService) error {
		_, err := s.Cron(r.cnf.Scheduler.HistoryExpression).
			SingletonMode().
			StartImmediately().
			Do(srv.Load, ctx)

		return err
	})
	if err != nil {
		return fmt.Errorf("history job scheduling error: %w", err)
	}

	err = r.container.Invoke(func(cnf *loader.Config) error {
		_, err = s.Cron(r.cnf.Scheduler.RefreshExpression).
			SingletonMode().
			Do(func() {
				//fmt.Printf("from refresh job\n")
			})

		return err
	})
	if err != nil {
		return fmt.Errorf("refresh job scheduling error: %w", err)
	}

	return nil
}

func (r *CompositionRoot) InitHttpServer() error {
	err := r.container.Invoke(func(h *Healtz) {
		http.Handle("/healtz", h)
	})
	if err != nil {
		return fmt.Errorf("health check handler error: %w", err)
	}
	return nil
}
