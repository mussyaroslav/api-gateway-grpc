package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorStrRecovery(c *gin.Context, recovered any) {
	if err, ok := recovered.(string); ok {
		// шаблон graceful обработки panic случая
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
