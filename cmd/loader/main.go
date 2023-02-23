package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/denis-gudim/economic-calendar/loader/loading"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func main() {
	root, err := NewCompositionRoot()
	if err != nil {
		err = fmt.Errorf("build application composition root failed: %w", err)
		processError(err)
	}
	defer root.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container := root.GetContainer()
	err = container.Invoke(func(s *loading.DictionariesLoaderService) error {
		return s.Load(ctx)
	})
	if err != nil {
		processError(err)
	}

	s := gocron.NewScheduler(time.UTC)
	if err = root.InitSchedule(ctx, s); err != nil {
		err = fmt.Errorf("init task scheduler failed: %w", err)
		processError(err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
	if err = root.InitHttpServer(); err != nil {
		err = fmt.Errorf("init http server failed: %w", err)
		processError(err)
	}

	s.StartAsync()

	log.Info("scheduler started...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			err = fmt.Errorf("listen http server failed: %w", err)
			processError(err)
		}
	}()

	<-ctx.Done()

	stop()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	s.Stop()

	log.Info("scheduler stoped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		err = fmt.Errorf("http server forced to shutdown: %w", err)
		processError(err)
	}
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
