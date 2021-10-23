package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"golang.org/x/xerrors"

	"github.com/denis-gudim/economic-calendar/api/app"
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

	root, err := app.NewCompositionRoot()

	if err != nil {
		err = xerrors.Errorf("build composition root failed: %w", err)
		processError(err)
	}

	defer root.Close()

	router := gin.Default()

	if err = root.InitRoutesV1(router); err != nil {
		err = xerrors.Errorf("init v1 api routes failed: %w", err)
		processError(err)
	}

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8080"); err != nil {
		err = xerrors.Errorf("run http server failed: %w", err)
		processError(err)
	}
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
