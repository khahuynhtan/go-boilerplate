package services

import (
	"myapp/api/author/entities"
	"myapp/api/author/repositories"
	"myapp/models"
	"myapp/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	authorRepo repositories.AuthorRepository
)

func InitService() {
	authorRepo = repositories.NewAuthorRepository()
	utils.Logger(utils.InfoLevel, "Author service initialized")
}

func GetListAuthors(c echo.Context) error {
	authors, err := authorRepo.GetList()
	if err != nil {
		return utils.Response(c, utils.StringPtr("OK"), authors, nil)
	}
	return utils.Response(c, utils.StringPtr("OK"), authors, nil)
}

func CreateAuthor(c echo.Context) error {
	validated_data := c.Get("validated").(*entities.CreateAuthorDto)
	author := &models.Author{
		Name:  validated_data.Name,
		Email: &validated_data.Email,
	}
	err := authorRepo.Create(author)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create author" + err.Error(),
		})
	}
	return utils.Response(c, utils.StringPtr("OK"), author, nil)
}
