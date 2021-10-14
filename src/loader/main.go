package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/denis-gudim/economic-calendar/loader/app"
	"github.com/denis-gudim/economic-calendar/loader/loading"
	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cnf := app.Config{}

	if err := cnf.Load(); err != nil {
		processError(err)
	}

	logger := log.StandardLogger()
	logger.SetLevel(cnf.Logging.Level)

	ds := loading.NewDictionariesLoaderService(cnf, logger)

	if err := ds.Load(); err != nil {
		processError(err)
	}

	s := gocron.NewScheduler(time.UTC)

	hs := loading.NewHistoryLoaderService(cnf, logger)

	_, err := s.Cron(cnf.Scheduler.HistoryExpression).
		SingletonMode().
		StartImmediately().
		Do(hs.Load)

	if err != nil {
		processError(err)
	}

	_, err = s.Cron(cnf.Scheduler.RefreshExpression).
		SingletonMode().
		Do(func() {
			fmt.Printf("from refresh job\n")
		})

	if err != nil {
		processError(err)
	}

	s.StartAsync()

	logger.Info("scheduler started...")

	<-ctx.Done()

	stop()

	s.Stop()

	logger.Info("scheduler stoped")
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
