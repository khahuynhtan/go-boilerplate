package utils

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Query(c echo.Context) *gorm.DB {
	db := c.Get("db").(*gorm.DB)
	return db
}
