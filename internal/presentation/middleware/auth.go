package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func BearerAuth(token string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			value := c.Request().Header.Get("Authorization")
			if strings.HasPrefix(value, "Bearer ") {
				value = strings.TrimPrefix(value, "Bearer ")
				if value != "" && value == token {
					return next(c)
				}
			}
			return c.JSON(http.StatusUnauthorized, struct {
				Message string `json:"message"`
			}{
				Message: "invalid token",
			})
		}
	}
}
