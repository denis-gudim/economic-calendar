package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/denis-gudim/economic-calendar/loader/loading"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

func main() {

	root, err := NewCompositionRoot()

	if err != nil {
		err = xerrors.Errorf("build composition root failed: %w", err)
		processError(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container := root.GetContainer()

	defer root.Close()

	err = container.Invoke(func(s *loading.DictionariesLoaderService) error {
		return s.Load(ctx)
	})

	if err != nil {
		processError(err)
	}

	s := gocron.NewScheduler(time.UTC)

	err = root.InitSchedule(ctx, s)

	if err != nil {
		err = xerrors.Errorf("init scheduler failed: %w", err)
		processError(err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	err = root.InitHttpServer()

	if err != nil {
		err = xerrors.Errorf("init http server failed: %w", err)
		processError(err)
	}

	s.StartAsync()

	log.Info("scheduler started...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			err = xerrors.Errorf("listen http server failed: %w", err)
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
		err = xerrors.Errorf("server forced to shutdown: %w", err)
		processError(err)
	}
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
