//go:build tools
// +build tools

package tools

import (
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/go-openapi/spec"
	_ "github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/joho/godotenv"
	_ "go.uber.org/zap"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)
