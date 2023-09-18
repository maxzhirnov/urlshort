package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/maxzhirnov/urlshort/internal/auth"
)

func TokenIssuerMiddleware(a *auth.Auth, l logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Смотрим есть ли уже jwt_token кука
		cookie, err := c.Cookie("jwt_token")
		if err == nil {
			// Проверяем токен на валидность
			_, err := a.ValidateToken(cookie)
			if err == nil {
				// если кука уже существует и токен валидный, то ничего не делаем
				return
			}
		}

		userID := a.GenerateUUID()
		jwtToken, err := a.GenerateToken(userID)
		if err != nil {
			l.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't generate token"})
			return
		}

		c.SetCookie("jwt_token", jwtToken, 3600*24*365, "/", "localhost", false, true)
		c.Next()
	}
}
