package todo

import (
	"ToDoProject/internal/models"
	"ToDoProject/tools"
	"database/sql"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (RouterTodo) createTask(c echo.Context) error {
	// Структура для тела запроса
	type RequestBody struct {
		TaskName        string  `json:"task_name" validate:"required"`
		TaskDescription *string `json:"task_description,omitempty"`
		DueDate         *string `json:"due_date,omitempty" validate:"datetime"`
		ParentTaskID    *uint64 `json:"parent_task_id,omitempty"`
	}

	var body RequestBody

	// Привязка данных запроса к структуре
	if err := c.Bind(&body); err != nil {
		return err
	}

	// Валидация данных запроса
	if err := c.Validate(&body); err != nil {
		return err
	}
	DueDate := sql.NullTime{Valid: false}
	if body.DueDate != nil {
		ParseDate, err := time.Parse("02.01.2006 15:04", *body.DueDate)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		DueDate.Time = ParseDate
		DueDate.Valid = true
	}

	db := tools.GetDBFromContext(c)

	user := c.Get("db_user").(*models.User)

	todo := models.Todo{
		TaskName:        body.TaskName,
		TaskDescription: body.TaskDescription,
		DueDate:         DueDate,
		ParentTaskID:    body.ParentTaskID,
		User:            user,
	}

	db.Create(&todo)
	return c.JSON(http.StatusOK, echo.Map{
		"todo": &todo,
	})
}
