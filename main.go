package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/andfxx27/my-cool-todo-list-restful-api/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type app struct {
	dbpool *pgxpool.Pool
}

type todo struct {
	Id          string     `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	CreatedDate time.Time  `db:"createdDate"`
	UpdatedDate *time.Time `db:"updatedDate"`
	DueDate     *time.Time `db:"dueDate"`
}

type saveTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
}

type response struct {
	Status  int
	Message string
	Result  interface{}
}

func (a app) saveTodo(c echo.Context) error {
	r := response{Status: http.StatusInternalServerError, Message: "Failed to save todo", Result: nil}

	body := new(saveTodoRequest)
	err := c.Bind(body)
	if err != nil {
		log.Err(errors.New("c.Bind error: " + err.Error()))
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
		log.Err(errors.New("pgxpool.Exec error: " + err.Error()))
		return c.JSON(http.StatusInternalServerError, r)
	}

	return c.JSON(http.StatusCreated, "saveTodo")
}

func (a app) getTodos(c echo.Context) error {
	r := response{Status: http.StatusInternalServerError, Message: "Failed to get todos", Result: nil}

	rows, _ := a.dbpool.Query(c.Request().Context(), "select * from todo")
	if rows.Err() != nil {
		log.Err(errors.New("pgxpool.Querry error: " + rows.Err().Error()))
		return c.JSON(http.StatusInternalServerError, r)
	}

	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[todo])
	if err != nil {
		log.Err(errors.New("pgx.CollectRows error: " + err.Error()))
		return c.JSON(http.StatusInternalServerError, r)
	}

	r = response{Status: http.StatusOK, Message: "Success get todos", Result: todos}

	return c.JSON(http.StatusOK, r)
}

func (a app) getTodo(c echo.Context) error {
	id := c.Param("id")

	rows, _ := a.dbpool.Query(c.Request().Context(), "select * from todo where id = $1", id)

	todo, _ := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[todo])

	r := response{Status: http.StatusOK, Message: "Success get todo", Result: todo}

	return c.JSON(http.StatusOK, r)
}

func updateTodo(c echo.Context) error {
	return c.JSON(http.StatusOK, "updateTodo")
}

func deleteTodo(c echo.Context) error {
	return c.JSON(http.StatusOK, "deleteTodo")
}

func main() {
	e := echo.New()

	dbpool := config.InitDatabaseConnection(e)

	e.Use(config.LoggerMiddleware())

	a := app{dbpool: dbpool}

	e.POST("/todos", a.saveTodo)
	e.GET("/todos", a.getTodos)
	e.GET("/todos/:id", a.getTodo)
	e.PUT("/todos/:id", updateTodo)
	e.DELETE("/todos/:id", deleteTodo)

	e.Logger.Fatal(e.Start(":1323"))
}
