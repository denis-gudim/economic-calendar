package main

import (
	"log"
	"os"

	"github.com/denis-gudim/economic-calendar/api/app"
	v1 "github.com/denis-gudim/economic-calendar/api/v1"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
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

	v1.InitRoutes(router, cnf, logger)

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	if err := router.Run(":8080"); err != nil {
		processError(err)
	}
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
