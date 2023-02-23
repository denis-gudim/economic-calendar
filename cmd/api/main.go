package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	//ginSwagger "github.com/swaggo/gin-swagger"
	// "github.com/swaggo/gin-swagger/swaggerFiles"
	ginprometheus "github.com/zsais/go-gin-prometheus"

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
	root, err := NewCompositionRoot()
	if err != nil {
		err = fmt.Errorf("build composition root failed: %w", err)
		processError(err)
	}
	defer root.Close()
	startHttpServer(root)
}

func startHttpServer(root *CompositionRoot) {
	router := gin.Default()

	if err := root.InitHttpServer(router); err != nil {
		err = fmt.Errorf("init http server failed: %w", err)
		processError(err)
	}

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			err = fmt.Errorf("listen http server failed: %w", err)
			processError(err)
		}
	}()

	<-ctx.Done()

	stop()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		err = fmt.Errorf("server forced to shutdown: %w", err)
		processError(err)
	}
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
