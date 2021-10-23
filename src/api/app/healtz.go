package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Healtz struct {
	db     *sqlx.DB
	checks map[string]func(ctx *gin.Context) (interface{}, error)
}

func NewHealtz(db *sqlx.DB) *Healtz {
	h := Healtz{db: db}

	h.checks = map[string]func(ctx *gin.Context) (interface{}, error){
		"db": h.checkDB,
	}

	return &h
}

func (h *Healtz) Handle(ctx *gin.Context) {
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
		res, err := v(ctx)

		if err != nil {
			failed++
		}

		out.Items[k] = res
	}

	out.Duration = time.Since(start).String()
	out.Failed = failed

	if failed > 0 {
		out.Status = "DOWN"
		ctx.JSON(http.StatusServiceUnavailable, out)
		return
	}

	ctx.JSON(http.StatusOK, out)
}

func (h *Healtz) checkDB(ctx *gin.Context) (interface{}, error) {
	out := struct {
		Status   string      `json:"status"`
		Result   interface{} `json:"result"`
		Duration string      `json:"duration"`
	}{
		Status: "UP",
	}

	count := 0
	start := time.Now()

	err := h.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM countries")

	out.Duration = time.Since(start).String()
	out.Result = count

	if err == nil && count == 0 {
		err = fmt.Errorf("validation result is empty")
	}

	if err != nil {
		out.Status = "DOWN"
		out.Result = err.Error()
	}

	return out, err
}
