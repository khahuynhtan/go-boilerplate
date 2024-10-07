package controllers

import (
	"myapp/api/author/services"

	"github.com/labstack/echo/v4"
)

type AuthController struct{}

func (AuthController AuthController) CreateAuthorHandler(c echo.Context) error {
	return services.CreateAuthor(c)
}

func (AuthController AuthController) GetListAuthorsHandler(c echo.Context) error {
	return services.GetListAuthors(c)
}
