package main

import (
	"net/http"
	"os"
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

	err = root.InitHttpServer()

	if err != nil {
		err = xerrors.Errorf("init http server failed: %w", err)
		processError(err)
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		err = xerrors.Errorf("start http server failed: %w", err)
		processError(err)
	}

	log.Info("scheduler stoped")
}

func processError(err error) {
	log.Fatal(err)
	os.Exit(2)
}
