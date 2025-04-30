package chef

import (
	"github.com/pkg/errors"
	"time"
)

const (
	timeoutRequestGRPC = 30 * time.Second
)

var (
	errInvalidID           = errors.New("Неверный параметр ID")
	errInvalidUUID         = errors.New("Неверный параметр UUID")
	errMissingRequiredUUID = errors.New("Пропущен обязательный параметр UUID")
	errSrvBadRequest       = errors.New("Ошибка при выполнении запроса к сервису")
	errSrvBadResponse      = errors.New("Ошибка при обработке данных от сервиса")
	errResNotFound         = errors.New("Запрошенный ресурс не найден")
	errInvalidArgument     = errors.New("Переданные аргументы неверны или недопустимы")
	errUnauthenticated     = errors.New("Запрос не прошел аутентификацию")
	errPermissionDenied    = errors.New("Нет прав доступа к запрошенному ресурсу")
	errDeadlineExceeded    = errors.New("Превышено время ожидания запроса")
)
