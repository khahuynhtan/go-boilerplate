package middlewares

import (
	"myapp/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Cors() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusOK)
			}
			utils.Logger(utils.InfoLevel, "Cors middleware executed")
			return next(c)
		}
	}
}
