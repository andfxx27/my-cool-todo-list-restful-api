package model

import "time"

type Todo struct {
	Id          string     `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	CreatedDate time.Time  `db:"createdDate"`
	UpdatedDate *time.Time `db:"updatedDate"`
	DueDate     *time.Time `db:"dueDate"`
}

type SaveTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"dueDate"`
}
