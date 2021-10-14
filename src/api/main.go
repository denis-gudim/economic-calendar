package main

import (
	"log"
	"os"

	"github.com/denis-gudim/economic-calendar/api/app"
	v1_handlers "github.com/denis-gudim/economic-calendar/api/v1/handlers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	cnf := app.Config{}

	if err := cnf.Load(); err != nil {
		processError(err)
	}

	logger, err := zap.NewProduction()

	if err != nil {
		processError(err)
	}

	defer logger.Sync()

	router := gin.Default()

	apiV1 := router.Group("/v1")

	v1_handlers.InitCountriesHandler(apiV1, cnf, logger)
	v1_handlers.InitScheduleHandler(apiV1, cnf, logger)

	if err := router.Run(":8081"); err != nil {
		processError(err)
	}
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
