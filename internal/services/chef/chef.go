package chef

import (
	"api-gateway-grpc/config"
	"api-gateway-grpc/internal/services/chef/middlewares"
	toAuthService "api-gateway-grpc/internal/services/grpc-client/to-auth-service"
	"api-gateway-grpc/pkg/logger"
	"context"
	"crypto/tls"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	log      *slog.Logger
	sendAuth *toAuthService.App
	address  string
	server   *http.Server
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func New(log *slog.Logger, sendAuth *toAuthService.App, cfg *config.Config) *Service {
	// инициализация http-сервера gin
	switch cfg.Env {
	case envDev, envProd:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	// инициализируем логгер в зависимости от уровня
	switch cfg.Env {
	case "local":
		router.Use(gin.Logger())
	default:
		router.Use(middlewares.StructuredLogger(log))
	}
	router.Use(gin.Recovery())
	router.Use(gin.CustomRecovery(middlewares.ErrorStrRecovery))

	var allowOrigins []string
	switch cfg.Env {
	case envDev, envProd:
		allowOrigins = append(allowOrigins, "https://")
		allowOrigins = append(allowOrigins, "https://")
	default:
		allowOrigins = append(allowOrigins, "*")
	}
	log.Info("CORS config", slog.String("allowOrigins", strings.Join(allowOrigins, ", ")))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,                                        // Разрешенные источники
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},   // Разрешенные методы
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Разрешенные заголовки
		ExposeHeaders:    []string{"Content-Length"},                          // Заголовки, которые могут быть доступны на клиенте
		AllowCredentials: true,                                                // Разрешить отправку кукисов
		MaxAge:           24 * time.Hour,                                      // Время кэширования preflight-запросов
	}))
	adr := cfg.Host + ":" + strconv.Itoa(cfg.Port)

	s := &Service{
		log:      log.With("service", "CHEF"),
		sendAuth: sendAuth,
		address:  adr,
		server: &http.Server{
			Addr:    adr,
			Handler: router,
		},
	}

	// установка всех маршрутов
	s.setRoutes(router)

	return s
}
func (s *Service) Run(cfg *config.Config) error {
	if cfg.Certs.Use {
		// Декодируем сертификаты
		cert, err := loadCerts(cfg.Certs.Crt, cfg.Certs.Key)
		if err != nil {
			s.log.Error("Failed to load certificates", logger.Err(err))
			return err
		}

		// Создаем TLS-конфигурацию
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		// Настраиваем сервер для использования TLS
		s.server.TLSConfig = tlsConfig

		// запуск gin https-сервера с TLS
		go func() {
			s.log.Info("Run HTTPS-server", slog.String("address", s.address))
			listener, err := tls.Listen("tcp", s.address, tlsConfig)
			if err != nil {
				s.log.Error("Failed to start TLS listener", logger.Err(err))
				return
			}

			if err := s.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.log.Error("Run https-server", logger.Err(err))
			}
		}()
	} else {
		// запуск gin http-сервера без TLS
		go func() {
			s.log.Info("Run HTTP-server", slog.String("address", s.address))
			if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.log.Error("Run http-server", logger.Err(err))
			}
		}()
	}

	return nil
}
func (s *Service) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("Shutdown http-server ", logger.Err(err))
	}

	// закрываем соединения с gRPC серверами
	s.sendAuth.Close()
}
