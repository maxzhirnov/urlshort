package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxzhirnov/urlshort/internal/auth"
)

const cookieExpireTime = 3600 * 24 * 365

func TokenIssuerMiddleware(a *auth.Auth, l logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем имеется ли кука jwt_token
		cookie, err := c.Request.Cookie("jwt_token")

		// Если куки нет или она не валидна
		if errors.Is(err, http.ErrNoCookie) || (err == nil && a.ValidateToken(cookie.Value) == "") {
			l.Error(err.Error())
			userID := a.GenerateUUID()
			jwtToken, err := a.GenerateToken(userID)
			if err != nil {
				l.Error("Failed to generate token: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't generate token"})
				return
			}
			c.SetCookie("jwt_token", jwtToken, cookieExpireTime, "/", "localhost", false, true)
			c.Set("jwt_token", jwtToken)
		} else if err != nil {
			// Вообще метод Request.Cookie() может вернуть только одну ошибку, но на всякий случай проверим
			l.Error("Unknown error: ", err)
		}
		c.Next()
	}
}
