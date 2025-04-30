package chef

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Service) ErrorBadRequest(c *gin.Context, err error) {
	data := gin.H{
		"title": "Bad request",
		"error": err.Error(),
	}
	c.JSON(http.StatusBadRequest, data)
}

func (s *Service) ErrorPageNotFound(c *gin.Context) {
	data := gin.H{
		"title": "Page not found",
	}
	c.Abort()
	c.JSON(http.StatusNotFound, data)
}

func (s *Service) ErrorServerError(c *gin.Context, err error) {
	data := gin.H{
		"title": "Server error",
		"error": err.Error(),
	}
	c.JSON(http.StatusInternalServerError, data)
}

func (s *Service) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"root": true,
	})
}

func (s *Service) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}
