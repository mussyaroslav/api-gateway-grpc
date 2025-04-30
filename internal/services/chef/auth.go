package chef

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	roleAdmin  = "admin"
	roleReader = "reader"
)

// Ошибки аутентификации
var (
	ErrNoAuthHeader      = errors.New("отсутствует заголовок авторизации")
	ErrInvalidAuthFormat = errors.New("неверный формат заголовка авторизации")
	ErrInvalidToken      = errors.New("недействительный токен")
)

// authMiddleware проверяет JWT токен и права доступа для указанной роли
func (s *Service) authMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrNoAuthHeader.Error()})
			c.Abort()
			return
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidAuthFormat.Error()})
			c.Abort()
			return
		}

		token := parts[1]

		// Проверяем токен через сервис аутентификации по gRPC
		tokenInfo, err := s.sendAuth.VerifyToken(c, token)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   fmt.Sprintf("ошибка проверки токена: %v", err),
				"message": "Пожалуйста, попробуйте перезайти или зарегистрироваться",
			})
			c.Abort()
			return
		}

		// Сохраняем данные пользователя в контексте
		c.Set("userId", tokenInfo.UserID)
		c.Set("userEmail", tokenInfo.Email)
		c.Set("userRoles", tokenInfo.Roles)

		// Если требуется роль reader, то достаточно быть аутентифицированным
		// (т.к. все зарегистрированные пользователи имеют эту роль)
		if role == roleReader {
			c.Next()
			return
		}

		// Для других ролей (например, admin) проверяем наличие требуемой роли
		hasRequiredRole := false
		for _, userRole := range tokenInfo.Roles {
			if userRole == role {
				hasRequiredRole = true
				break
			}
		}

		if !hasRequiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   fmt.Sprintf("у вас нет прав к этому функционалу (требуется роль: %s)", role),
				"message": "Пожалуйста, обратитесь к администратору для получения необходимых прав",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Вспомогательные функции для получения данных пользователя из контекста

// GetUserID возвращает ID пользователя из контекста
func GetUserID(c *gin.Context) (string, bool) {
	userId, exists := c.Get("userId")
	if !exists {
		return "", false
	}
	userIdStr, ok := userId.(string)
	return userIdStr, ok
}

// GetUserEmail возвращает email пользователя из контекста
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("userEmail")
	if !exists {
		return "", false
	}
	emailStr, ok := email.(string)
	return emailStr, ok
}

// HasRole проверяет, имеет ли пользователь указанную роль
func HasRole(c *gin.Context, role string) bool {
	rolesInterface, exists := c.Get("userRoles")
	if !exists {
		return false
	}

	roles, ok := rolesInterface.([]string)
	if !ok {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

// authReader возвращает middleware для проверки базовой аутентификации
// (все зарегистрированные пользователи имеют роль reader)
func (s *Service) authReader() gin.HandlerFunc {
	return s.authMiddleware(roleReader)
}

// authAdmin возвращает middleware для проверки роли администратора
func (s *Service) authAdmin() gin.HandlerFunc {
	return s.authMiddleware(roleAdmin)
}
