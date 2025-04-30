package to_auth_service

import (
	"api-gateway-grpc/config"
	apiAuthService "api-gateway-grpc/generate/auth-service"
	clientGRPC "api-gateway-grpc/pkg/grpc/client"
	"api-gateway-grpc/pkg/logger"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
	"log/slog"
	"math/rand"
	"time"
)

type App struct {
	log        *slog.Logger
	grpcConn   *grpc.ClientConn
	grpcClient apiAuthService.AuthServiceClient
	cfgClient  config.ClientGRPC
}

const nameTo = "AUTH-SERVICE"

func MustNew(log *slog.Logger, cfg config.ClientGRPC) *App {
	const proc = "SENDER to " + nameTo

	grpcConn, err := clientGRPC.New(&cfg)
	if err != nil {
		panic("create new client grpc to " + nameTo + ": " + err.Error())
	}

	return &App{
		log:        log.With(slog.String("proc", proc)),
		grpcClient: apiAuthService.NewAuthServiceClient(grpcConn),
		grpcConn:   grpcConn,
		cfgClient:  cfg,
	}
}

func (a *App) Close() {
	if err := a.grpcConn.Close(); err != nil {
		a.log.Error("close grpc connection to "+nameTo, logger.Err(err))
	} else {
		a.log.Debug("grpc connection to " + nameTo + " closed")
	}
}

func (a *App) checkConnectionState() error {
	// проверим, находится ли соединение в состоянии сбоя
	if a.grpcConn.GetState() != connectivity.TransientFailure {
		return nil
	}
	// соединение в состоянии сбоя(TransientFailure), ждем его изменения(EXPERIMENTAL!)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutWaitChangeState)
	defer cancel()
	if !a.grpcConn.WaitForStateChange(ctx, connectivity.TransientFailure) {
		return status.Errorf(codes.Unavailable, "connection recovery timed out")
	}

	return nil
}

func (a *App) executeWithRetries(ctx context.Context, fn func(_ context.Context) (interface{}, error)) (interface{}, error) {
	backoff := minBackoff
	for retry := 0; retry < maxRetries; retry++ {
		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}

		if e, ok := status.FromError(err); ok && (e.Code() == codes.Unavailable || e.Code() == codes.DeadlineExceeded) {
			delay := float64(backoff) * (1 + jitter*(rand.Float64()*2-1))
			time.Sleep(time.Duration(delay))
			backoff *= 2

			//a.log.Debug("retries",
			//	logger.Err(err),
			//	slog.Int("retry", retry),
			//	slog.Duration("delay", time.Duration(delay)),
			//)

			if err := a.checkConnectionState(); err != nil {
				return nil, err
			}

			continue
		}

		return nil, err
	}

	return nil, errors.New("unsuccessful series of request tries")
}
