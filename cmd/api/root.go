package main

import (
	"fmt"

	"github.com/denis-gudim/economic-calendar/api"
	v1_controllers "github.com/denis-gudim/economic-calendar/api/v1/controllers"
	v1_data "github.com/denis-gudim/economic-calendar/api/v1/data"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type CompositionRoot struct {
	logger    *zap.Logger
	db        *sqlx.DB
	container *dig.Container
}

func NewCompositionRoot() (*CompositionRoot, error) {
	container := dig.New()
	cnf := api.Config{}

	if err := cnf.Load(); err != nil {
		return nil, fmt.Errorf("load application config error: %w", err)
	}

	logger, err := zap.NewProduction()

	if err != nil {
		return nil, fmt.Errorf("logger initialization error: %w", err)
	}

	db, err := sqlx.Connect("postgres", cnf.DB.ConnectionString)

	if err != nil {
		return nil, fmt.Errorf("connect to db error: %w", err)
	}

	container.Provide(func() *api.Config {
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

	container.Provide(NewHealtz)

	return &CompositionRoot{
		db:        db,
		logger:    logger,
		container: container,
	}, nil
}

func (r *CompositionRoot) InitHttpServer(gin *gin.Engine) error {
	v1 := gin.Group("v1")

	err := r.container.Invoke(func(c *v1_controllers.CountriesController) {
		g := v1.Group("countries")

		g.GET("", c.GetByLanguage)
	})

	if err != nil {
		return fmt.Errorf("countries controller init error: %w", err)
	}

	err = r.container.Invoke(func(c *v1_controllers.EventsController) {
		g := v1.Group("events")

		g.GET("", c.GetEventsSchedule)
		g.GET(":eventId", c.GetEventDetails)
		g.GET(":eventId/history", c.GetEventHistory)
	})

	if err != nil {
		return fmt.Errorf("events controller init error: %w", err)
	}

	err = r.container.Invoke(func(c *Healtz) {
		gin.GET("/healtz", c.Handle)
	})

	if err != nil {
		return fmt.Errorf("healtz controller init error: %w", err)
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
