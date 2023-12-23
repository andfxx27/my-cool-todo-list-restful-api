package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func InitDatabaseConnection(e *echo.Echo) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:superuser@localhost:5432/db-my-cool-todo-list-restful-api")
	if err != nil {
		e.Logger.Fatal(err.Error())
	}

	return dbpool
}
