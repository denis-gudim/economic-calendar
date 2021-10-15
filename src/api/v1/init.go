package v1

import (
	"github.com/denis-gudim/economic-calendar/api/app"
	handlers "github.com/denis-gudim/economic-calendar/api/v1/handlers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRoutes(gin *gin.Engine, cnf app.Config, logger *zap.Logger) {

	apiGroup := gin.Group("v1")

	// countries
	ch := handlers.NewCountriesHandler(cnf, logger)
	apiGroup.GET("countries", ch.Get)

	// events
	eh := handlers.NewEventsHandler(cnf, logger)
	apiGroup.GET("events", eh.GetEventsSchdule)
	apiGroup.GET("events/:id", eh.GetEventDetails)
	apiGroup.GET("events/:id/history", eh.GetEventHistory)
}
