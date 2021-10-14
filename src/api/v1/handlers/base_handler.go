package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type baseHandler struct {
}

func (h *baseHandler) writeServerError(c *gin.Context) {
	c.String(http.StatusInternalServerError, "500 Internal Server Error")
}

func (h *baseHandler) writeBadRequest(c *gin.Context, format string, values ...interface{}) {
	c.String(http.StatusBadRequest, "400 Bad Request: "+format, values...)
}
