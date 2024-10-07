package server

import (
	"myapp/middlewares"
	"myapp/routes"
	"myapp/scripts"
	"myapp/utils"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerType struct {
	Echo *echo.Echo
}

func NewServer() *ServerType {
	return &ServerType{
		Echo: echo.New(),
	}
}

func (s *ServerType) Start() error {
	err := godotenv.Load()
	if err != nil {
		utils.Logger(utils.FatalLevel, "Error loading .env file")
		return err
	}

	s.Echo.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	dbConnector := utils.NewConnector(
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	db := dbConnector.Connect()

	if db == nil {
		utils.Logger(utils.FatalLevel, "Failed to connect to DB")
		return nil
	}

	// Middlewares
	s.Echo.Use(middleware.Logger())
	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middlewares.Cors())

	var routingInstance routes.Routing
	endpoints := routingInstance.RegisterRoutes()

	for _, ep := range endpoints {
		s.Echo.Add(ep.Method, ep.Path, echo.HandlerFunc(ep.Handler), ep.Middlewares...)
	}

	scripts.GenerateOpenAPI(endpoints, "API Documentation", "1.0.0")

	// init services
	utils.Logger(utils.InfoLevel, "Initializing services...")
	scripts.InitServices()
	return nil
}
