package to_auth_service

import (
	apiAuthService "api-gateway-grpc/generate/auth-service"
	"api-gateway-grpc/internal/services/models"
	"context"
	"github.com/pkg/errors"
)

func (a *App) Ping(ctx context.Context) (bool, error) {
	result, err := a.executeWithRetries(ctx, func(_ context.Context) (interface{}, error) {
		return a.grpcClient.Ping(ctx, &apiAuthService.PingRequest{})
	})
	if err != nil {
		return false, errors.Wrap(err, "request Ping()")
	}

	return result.(*apiAuthService.PingResponse).GetOk(), nil
}

func (a *App) Register(ctx context.Context, request *models.AuthRequest) (*models.AuthResponse, error) {
	result, err := a.executeWithRetries(ctx, func(_ context.Context) (interface{}, error) {
		return a.grpcClient.Register(ctx, &apiAuthService.RegisterRequest{
			Email:    request.Email,
			Password: request.Password,
		})
	})
	if err != nil {
		return nil, errors.Wrap(err, "request Register()")
	}

	// Приводим результат к нужному типу
	resp, ok := result.(*apiAuthService.LoginResponse)
	if !ok {
		return nil, errors.New("unexpected response type from Register()")
	}

	registerResponse := &models.AuthResponse{
		JWTToken: resp.JwtToken,
	}

	return registerResponse, nil
}

func (a *App) Login(ctx context.Context, request *models.AuthRequest) (*models.AuthResponse, error) {
	result, err := a.executeWithRetries(ctx, func(_ context.Context) (interface{}, error) {
		return a.grpcClient.Login(ctx, &apiAuthService.LoginRequest{
			Email:    request.Email,
			Password: request.Password,
		})
	})
	if err != nil {
		return nil, errors.Wrap(err, "request Login()")
	}

	// Приводим результат к нужному типу
	resp, ok := result.(*apiAuthService.LoginResponse)
	if !ok {
		return nil, errors.New("unexpected response type from Login()")
	}

	loginResponse := &models.AuthResponse{
		JWTToken: resp.JwtToken,
	}

	return loginResponse, nil
}

func (a *App) VerifyToken(ctx context.Context, token string) (*models.TokenInfo, error) {
	result, err := a.executeWithRetries(ctx, func(_ context.Context) (interface{}, error) {
		return a.grpcClient.VerifyToken(ctx, &apiAuthService.VerifyTokenRequest{Token: token})
	})
	if err != nil {
		return nil, errors.Wrap(err, "request VerifyToken()")
	}

	// Приводим результат к нужному типу
	resp, ok := result.(*apiAuthService.VerifyTokenResponse)
	if !ok {
		return nil, errors.New("unexpected response type from VerifyToken()")
	}

	// Проверяем валидность токена
	if !resp.Valid {
		errorMsg := "недействительный токен"
		if resp.Error != nil && resp.Error.Message != "" {
			errorMsg = resp.Error.Message
		}
		return nil, errors.New(errorMsg)
	}

	// Создаем и возвращаем информацию о jwt токене
	tokenInfo := &models.TokenInfo{
		UserID:  resp.UserId,
		Email:   resp.Email,
		Roles:   resp.Roles,
		IsValid: resp.Valid,
	}

	return tokenInfo, nil
}
