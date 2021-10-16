package main

import (
	"log"
	"os"

	"github.com/denis-gudim/economic-calendar/api/app"
	v1 "github.com/denis-gudim/economic-calendar/api/v1"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"

	_ "github.com/denis-gudim/economic-calendar/api/docs"
)

// @title Economic Calendar Example API
// @version 1.0
// @description This is a sample economic calendar api

// @contact.name API Support
// @contact.url https://github.com/denis-gudim/economic-calendar/issues
// @contact.email denis.gudim@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1/
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8080"); err != nil {
		processError(err)
	}
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
