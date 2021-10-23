package app

import (
	v1_controllers "github.com/denis-gudim/economic-calendar/api/v1/controllers"
	v1_data "github.com/denis-gudim/economic-calendar/api/v1/data"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type CompositionRoot struct {
	logger    *zap.Logger
	db        *sqlx.DB
	container *dig.Container
}

func NewCompositionRoot() (*CompositionRoot, error) {
	container := dig.New()
	cnf := Config{}

	if err := cnf.Load(); err != nil {
		return nil, xerrors.Errorf("load config error: %w", err)
	}

	logger, err := zap.NewProduction()

	if err != nil {
		return nil, xerrors.Errorf("init logger error: %w", err)
	}

	db, err := sqlx.Connect("postgres", cnf.DB.ConnectionString)

	if err != nil {
		return nil, xerrors.Errorf("connect to db error: %w", err)
	}

	container.Provide(func() *Config {
		return &cnf
	})

	container.Provide(func() *zap.Logger {
		return logger
	})

	container.Provide(func() *sqlx.DB {
		return db
	})

	container.Provide(func(db *sqlx.DB) v1_controllers.CountriesDataReciver {
		return v1_data.NewCountriesRepository(db)
	})
	container.Provide(func(db *sqlx.DB) v1_controllers.EventsDataReciver {
		return v1_data.NewEventsRepository(db)
	})

	container.Provide(v1_controllers.NewCountriesController)
	container.Provide(v1_controllers.NewEventsController)

	return &CompositionRoot{
		db:        db,
		logger:    logger,
		container: container,
	}, nil
}

func (r *CompositionRoot) InitRoutesV1(gin *gin.Engine) error {
	v1 := gin.Group("v1")

	err := r.container.Invoke(func(c *v1_controllers.CountriesController) {
		g := v1.Group("countries")

		g.GET("", c.GetByLanguage)
	})

	if err != nil {
		return xerrors.Errorf("countries countroller error: %w", err)
	}

	err = r.container.Invoke(func(c *v1_controllers.EventsController) {
		g := v1.Group("events")

		g.GET("", c.GetEventsSchedule)
		g.GET(":eventId", c.GetEventDetails)
		g.GET(":eventId/history", c.GetEventHistory)
	})

	if err != nil {
		return xerrors.Errorf("events countroller error: %w", err)
	}

	return nil
}

func (r *CompositionRoot) Close() {
	defer func() {
		if r.logger != nil {
			r.logger.Sync()
		}
	}()

	if r.db != nil {
		r.db.Close()
	}
}
