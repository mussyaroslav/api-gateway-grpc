package chef

import (
	"api-gateway-grpc/config"
	"api-gateway-grpc/internal/services/chef"
	toAuthService "api-gateway-grpc/internal/services/grpc-client/to-auth-service"
	"log/slog"
)

type App struct {
	log     *slog.Logger
	Service *chef.Service
}

const serviceName = "CHEF"

// New create new service Chef
func New(
	log *slog.Logger,
	sendAuth *toAuthService.App,
	cfg *config.Config,
) *App {
	// GIN-сервер как сервис Chef
	srv := chef.New(log, sendAuth, cfg)

	return &App{
		log:     log,
		Service: srv,
	}
}

// MustRun runs service Chef and panics if any error occurs
func (a *App) MustRun(cfg *config.Config) {
	if err := a.Run(cfg); err != nil {
		panic(err)
	}
}

// Run runs service Chef
func (a *App) Run(cfg *config.Config) error {
	a.log.Info("start " + serviceName + " service...")
	return a.Service.Run(cfg)
}

// Stop stops RED service
func (a *App) Stop() {
	a.Service.Stop()
	a.log.Info(serviceName + " service stopped")
}
