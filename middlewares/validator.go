package middlewares

import (
	"myapp/utils"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func ValidateRequestMiddleware(model interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create a new instance of the model for each request
			model := reflect.New(reflect.TypeOf(model).Elem()).Interface()
			// Bind and validate the request
			if err := c.Bind(model); err != nil {
				return utils.Response(c, utils.StringPtr("Invalid request format"), nil, utils.IntPtr(http.StatusBadRequest))
			}

			// Validate using custom validator
			if err := c.Validate(model); err != nil {
				return utils.Response(c, utils.StringPtr("Validation failed"), err.Error(), utils.IntPtr(http.StatusUnprocessableEntity))
			}
			// assign validated data into the context
			c.Set("validated", model)
			return next(c)
		}
	}
}
