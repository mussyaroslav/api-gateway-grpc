package chef

import (
	"github.com/gin-gonic/gin"
)

func (s *Service) setRoutes(r *gin.Engine) {
	r.NoRoute(s.ErrorPageNotFound)
	r.NoMethod(s.ErrorPageNotFound)

	// общие запросы без валидации
	r.GET("/health", s.Health)

	api := r.Group("/api")
	{
		// Публичные маршруты (без аутентификации)
		public := api.Group("")
		{
			public.GET("/some-public-endpoint")
		}

		// Маршруты требующие аутентификации
		authenticated := api.Group("")
		authenticated.Use(s.authReader())
		{
			authenticated.GET("/some-private-endpoint")
		}
	}
}
