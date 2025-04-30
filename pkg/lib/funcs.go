package lib

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
)

// MaskedText возвращает замаскированную точками строку, оставив слева и справа оригинал в n символов
func MaskedText(msg string, n int) string {
	c := int(math.Floor(float64(len(msg)) / 3))
	if c > n {
		c = n
	}
	if c <= 0 {
		c = 1
	}
	return msg[0:c] + "..." + msg[len(msg)-c:]
}

// IsDeadlineExceeded проверяет, является ли ошибка 'err' ошибкой 'context.DeadlineExceeded'
func IsDeadlineExceeded(err error) bool {
	// Проверяем, была ли ошибка обернута в grpc.Status
	s, ok := status.FromError(err)
	if ok {
		return s.Code() == codes.DeadlineExceeded
	}

	// Если ошибка не является grpc.Status, проверяем явно тип ошибки
	return errors.Is(err, context.DeadlineExceeded)
}
