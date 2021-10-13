package main

import (
	"economic-calendar/loader/app"
	"economic-calendar/loader/loading"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
)

func main() {

	cnf := app.Config{}

	if err := cnf.Load(); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	logger := log.StandardLogger()
	logger.SetLevel(cnf.Logging.Level)

	ds := loading.NewDictionariesLoaderService(cnf, logger)

	if err := ds.Load(); err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	s := gocron.NewScheduler(time.UTC)

	hs := loading.NewHistoryLoaderService(cnf, logger)

	s.Cron(cnf.Scheduler.HistoryExpression).
		SingletonMode().
		StartImmediately().
		Do(hs.Load)

	s.Cron(cnf.Scheduler.RefreshExpression).
		SingletonMode().
		Do(func() {
			fmt.Printf("from refresh job\n")
		})

	s.StartAsync()

	logger.Info("scheduler started...")

	<-sigs

	s.Stop()

	logger.Info("scheduler stoped")
}
