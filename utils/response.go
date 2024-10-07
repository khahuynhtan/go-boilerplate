package utils

import "github.com/labstack/echo/v4"

type ResponseType struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	StatusCode int         `json:"status_code"`
}

func Response(c echo.Context, message *string, data interface{}, statusCode *int) error {
	// Set default values if nil
	defaultMessage := "Success"
	defaultStatusCode := 200

	if message == nil {
		message = &defaultMessage
	}
	if statusCode == nil {
		statusCode = &defaultStatusCode
	}

	return c.JSON(*statusCode, &ResponseType{
		Message:    *message,
		Data:       data,
		StatusCode: *statusCode,
	})
}
