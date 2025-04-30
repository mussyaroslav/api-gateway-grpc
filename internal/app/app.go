package app

import (
	"api-gateway-grpc/config"
	chefApp "api-gateway-grpc/internal/app/chef"
	toAuthService "api-gateway-grpc/internal/services/grpc-client/to-auth-service"
	"api-gateway-grpc/pkg/logger"
	"context"
	"log/slog"
	"time"
)

type App struct {
	log             *slog.Logger
	SendAuthService *toAuthService.App
	Chef            *chefApp.App
}

const (
	pingTimeout = 5 * time.Second
)

func New(log *slog.Logger, cfg *config.Config) *App {
	appSendAuthService := toAuthService.MustNew(log, cfg.AuthService)

	// GIN-сервер как application
	chef := chefApp.New(log, appSendAuthService, cfg)

	return &App{
		log:             log,
		SendAuthService: appSendAuthService,
		Chef:            chef,
	}
}

func (a *App) MustRun(cfg *config.Config) {
	var err error

	ctxPARH, cancelPARH := context.WithTimeout(context.Background(), pingTimeout)
	defer cancelPARH()
	_, err = a.SendAuthService.Ping(ctxPARH)
	if err != nil {
		a.log.Warn("ping of AUTH-SERVICE", logger.Err(err))
	} else {
		a.log.Debug("success ping of AUTH-SERVICE")
	}

	a.Chef.MustRun(cfg)
}

func (a *App) Stop() {
	a.Chef.Stop()
}
