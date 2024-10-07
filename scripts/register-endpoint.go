package scripts

import (
	"myapp/types"

	"github.com/labstack/echo/v4"
)

func RegisterEndpoint(path string, method string, summary string, operationID string, extendSchema types.ExtendSchema, handler types.APIHandler, endpoints *[]types.Endpoint, middlewares ...echo.MiddlewareFunc) {
	*endpoints = append(*endpoints, types.Endpoint{Path: path, Method: method, Summary: summary, OperationID: operationID, Handler: handler, Middlewares: middlewares, ExtendSchema: extendSchema})
}
