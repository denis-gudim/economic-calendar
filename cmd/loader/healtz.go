package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Healtz struct {
	db     *sql.DB
	logger *logrus.Logger
	checks map[string]func(req *http.Request) (interface{}, error)
}

func NewHealtz(db *sql.DB, logger *logrus.Logger) *Healtz {
	h := Healtz{db: db, logger: logger}

	h.checks = map[string]func(req *http.Request) (interface{}, error){
		"db": h.checkDB,
	}

	return &h
}

func (h *Healtz) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	out := struct {
		Status   string                 `json:"status"`
		Duration string                 `json:"duration"`
		Failed   int                    `json:"failed"`
		Items    map[string]interface{} `json:"items"`
	}{
		Status: "UP",
		Items:  make(map[string]interface{}),
	}

	failed := 0
	start := time.Now()

	for k, v := range h.checks {
		res, err := v(req)

		if err != nil {
			failed++
			h.logger.Error(err)
		}

		out.Items[k] = res
	}

	out.Duration = time.Since(start).String()
	out.Failed = failed

	if failed > 0 {
		out.Status = "DOWN"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
		h.logger.WithField("duration", out.Duration).Info("healthcheck complete successfully")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Healtz) checkDB(req *http.Request) (interface{}, error) {
	out := struct {
		Status   string      `json:"status"`
		Result   interface{} `json:"result"`
		Duration string      `json:"duration"`
	}{
		Status: "UP",
	}

	count := 0
	start := time.Now()
	row := h.db.QueryRowContext(req.Context(), "SELECT COUNT(*) FROM countries")

	err := row.Scan(&count)
	if err != nil {
		out.Status = "DOWN"
		out.Result = err.Error()
		err = fmt.Errorf("database healthcheck err: %w", err)
	} else if count == 0 {
		err = fmt.Errorf("validation result is empty")
	}

	out.Duration = time.Since(start).String()
	out.Result = count

	return out, err
}
