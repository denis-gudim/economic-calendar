package main

import (
	"context"
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

	container := root.GetContainer()

	defer root.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	err = container.Invoke(func(s *loading.DictionariesLoaderService) error {
		return s.Load()
	})

	if err != nil {
		processError(err)
	}

	s := gocron.NewScheduler(time.UTC)

	err = root.InitSchedule(s)

	if err != nil {
		err = xerrors.Errorf("init scheduler failed: %w", err)
		processError(err)
	}

	s.StartAsync()

	log.Info("scheduler started...")

	<-ctx.Done()

	stop()

	s.Stop()

	log.Info("scheduler stoped")
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
