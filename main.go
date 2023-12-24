package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/andfxx27/my-cool-todo-list-restful-api/config"
	"github.com/andfxx27/my-cool-todo-list-restful-api/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type app struct {
	dbpool *pgxpool.Pool
}

func (a app) saveTodo(c echo.Context) error {
	r := model.HTTPResponse{Status: http.StatusInternalServerError, Message: "Failed to save todo", Result: nil}

	body := new(model.SaveTodoRequest)
	err := c.Bind(body)
	if err != nil {
		log.Err(errors.New("c.Bind error: " + err.Error())).Msg("saveTodo error")
		r.Status = http.StatusBadRequest
		return c.JSON(http.StatusBadRequest, r)
	}

	// TODO Check if there exists todo with same title

	dueDate, _ := time.Parse(time.RFC3339, body.DueDate)

	_, err = a.dbpool.Exec(
		c.Request().Context(),
		"insert into todo (id, title, description, status, dueDate) values ($1,$2,$3,$4,$5)",
		uuid.New().String(), body.Title, body.Description, "TODO", dueDate,
	)
	if err != nil {
		log.Err(errors.New("pgxpool.Exec error: " + err.Error())).Msg("saveTodo error")
		return c.JSON(http.StatusInternalServerError, r)
	}

	r = model.HTTPResponse{Status: http.StatusCreated, Message: "Success save todo", Result: nil}

	return c.JSON(http.StatusCreated, r)
}

func (a app) getTodos(c echo.Context) error {
	r := model.HTTPResponse{Status: http.StatusInternalServerError, Message: "Failed to get todos", Result: nil}

	rows, _ := a.dbpool.Query(c.Request().Context(), "select * from todo")
	if rows.Err() != nil {
		log.Err(errors.New("pgxpool.Querry error: " + rows.Err().Error())).Msg("getTodos error")
		return c.JSON(http.StatusInternalServerError, r)
	}
	defer rows.Close()

	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Todo])
	if err != nil {
		log.Err(errors.New("pgx.CollectRows error: " + err.Error())).Msg("getTodos error")
		return c.JSON(http.StatusInternalServerError, r)
	}

	r = model.HTTPResponse{Status: http.StatusOK, Message: "Success get todos", Result: todos}

	return c.JSON(http.StatusOK, r)
}

func (a app) getTodo(c echo.Context) error {
	r := model.HTTPResponse{Status: http.StatusInternalServerError, Message: "Failed to get todo", Result: nil}

	id := c.Param("id")

	rows, _ := a.dbpool.Query(c.Request().Context(), "select * from todo where id = $1", id)
	if rows.Err() != nil {
		log.Err(errors.New("pgxpool.Querry error: " + rows.Err().Error())).Msg("getTodo error")
		return c.JSON(http.StatusInternalServerError, r)
	}
	defer rows.Close()

	todo, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Todo])
	if err != nil {
		log.Err(errors.New("pgx.CollectExactlyOneRow error: " + err.Error())).Msg("getTodo error")
		return c.JSON(http.StatusInternalServerError, r)
	}

	r = model.HTTPResponse{Status: http.StatusOK, Message: "Success get todo", Result: todo}

	return c.JSON(http.StatusOK, r)
}

func (a app) updateTodo(c echo.Context) error {
	r := model.HTTPResponse{Status: http.StatusInternalServerError, Message: "Failed to update todo", Result: nil}

	id := c.Param("id")

	body := new(model.UpdateTodoRequest)
	err := c.Bind(body)
	if err != nil {
		log.Err(errors.New("c.Bind error: " + err.Error())).Msg("updateTodo error")
		r.Status = http.StatusBadRequest
		return c.JSON(http.StatusBadRequest, r)
	}

	_, err = a.dbpool.Exec(
		c.Request().Context(),
		"update todo set title = $1, description = $2, status = $3, dueDate = $4 where id = $5",
		body.Title, body.Description, body.Status, body.DueDate, id,
	)
	if err != nil {
		log.Err(errors.New("pgxpool.Exec error: " + err.Error())).Msg("updateTodo error")
		return c.JSON(http.StatusInternalServerError, r)
	}

	r = model.HTTPResponse{Status: http.StatusOK, Message: "Success update todo", Result: nil}

	return c.JSON(http.StatusOK, r)
}

func (a app) deleteTodo(c echo.Context) error {
	r := model.HTTPResponse{Status: http.StatusInternalServerError, Message: "Failed to delete todo", Result: nil}

	id := c.Param("id")

	_, err := a.dbpool.Exec(c.Request().Context(), "delete from todo where id = $1", id)
	if err != nil {
		log.Err(errors.New("pgxpool.Exec error: " + err.Error())).Msg("deleteTodo error")
		return c.JSON(http.StatusInternalServerError, r)
	}

	r = model.HTTPResponse{Status: http.StatusOK, Message: "Success delete todo", Result: nil}

	return c.JSON(http.StatusOK, r)
}

func main() {
	e := echo.New()

	dbpool := config.InitDatabaseConnection(e)

	e.Use(config.LoggerMiddleware())

	a := app{dbpool: dbpool}

	e.POST("/todos", a.saveTodo)
	e.GET("/todos", a.getTodos)
	e.GET("/todos/:id", a.getTodo)
	e.PUT("/todos/:id", a.updateTodo)
	e.DELETE("/todos/:id", a.deleteTodo)

	e.Logger.Fatal(e.Start(":1323"))
}
