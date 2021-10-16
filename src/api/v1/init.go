package v1

import (
	"github.com/denis-gudim/economic-calendar/api/app"
	controllers "github.com/denis-gudim/economic-calendar/api/v1/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRoutes(gin *gin.Engine, cnf app.Config, logger *zap.Logger) {

	c := controllers.NewCountriesController(cnf, logger)
	e := controllers.NewEventsController(cnf, logger)

	v1 := gin.Group("v1")
	{
		countires := v1.Group("countries")
		{
			countires.GET("", c.GetByLanguage)
		}
		events := v1.Group("events")
		{
			events.GET("", e.GetEventsSchedule)
			events.GET(":eventId", e.GetEventDetails)
			events.GET(":eventId/history", e.GetEventHistory)
		}
	}
}
