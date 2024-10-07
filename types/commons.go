package types

import "github.com/labstack/echo/v4"

type APIHandler func(c echo.Context) error

type Parameter struct {
	Name        string
	Type        string
	Location    string
	Required    bool
	Description string
	Examples    interface{}
	Default     interface{}
}

type ExtendSchema struct {
	IsAuth      *bool
	Parameters  []Parameter
	Request     interface{}
	Responses   map[int]interface{}
	Description string
	Tag         string
}

type Endpoint struct {
	Middlewares  []echo.MiddlewareFunc
	ExtendSchema ExtendSchema
	Path         string
	Method       string
	Summary      string
	OperationID  string
	Handler      APIHandler
	Description  *string
}
