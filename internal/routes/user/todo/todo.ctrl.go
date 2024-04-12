package todo

import (
	"ToDoProject/internal/models"
	"ToDoProject/tools"
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
	var DueDate *time.Time
	if body.DueDate != nil {
		ParseDate, err := time.Parse("02.01.2006 15:04", *body.DueDate)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		DueDate = &ParseDate
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

func (RouterTodo) updateTask(c echo.Context) error {
	type RequestBody struct {
		TaskID          uint64  `json:"task_id" validate:"required"`
		TaskName        *string `json:"task_name,omitempty"`
		TaskDescription *string `json:"task_description,omitempty"`
		DueDate         *string `json:"due_date,omitempty" validate:"omitempty,datetime"`
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

	var DueDate *time.Time
	if body.DueDate != nil {
		ParseDate, err := time.Parse("02.01.2006 15:04", *body.DueDate)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		DueDate = &ParseDate
	}

	db := tools.GetDBFromContext(c)

	user := c.Get("db_user").(*models.User)

	var todo models.Todo
	if result := db.Where("id = ? AND user_id = ?", body.TaskID, user.ID).Find(&todo); result.Error != nil || result.RowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	if body.TaskName != nil {
		todo.TaskName = *body.TaskName
	}
	if body.TaskDescription != nil {
		todo.TaskDescription = body.TaskDescription
	}
	if DueDate != nil {
		todo.DueDate = DueDate
	}

	db.Save(&todo)

	return c.JSON(http.StatusOK, echo.Map{
		"todo": &todo,
	})
}

func (RouterTodo) switchStatusTask(c echo.Context) error {
	type RequestBody struct {
		TaskID uint64 `json:"task_id" validate:"required"`
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

	db := tools.GetDBFromContext(c)

	user := c.Get("db_user").(*models.User)

	var todo models.Todo
	if result := db.Where("id = ? AND user_id = ?", body.TaskID, user.ID).Find(&todo); result.Error != nil || result.RowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	todo.Completed = !todo.Completed
	db.Save(&todo)
	return c.JSON(http.StatusOK, echo.Map{
		"todo": &todo,
	})
}

func (RouterTodo) deleteTask(c echo.Context) error {
	type RequestBody struct {
		TaskID uint64 `json:"task_id" validate:"required"`
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

	db := tools.GetDBFromContext(c)

	user := c.Get("db_user").(*models.User)

	var todo models.Todo
	if result := db.Where("id = ? AND user_id = ?", body.TaskID, user.ID).Find(&todo); result.Error != nil || result.RowsAffected == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	db.Delete(&todo)
	return c.NoContent(http.StatusNoContent)
}

func (RouterTodo) getTasks(c echo.Context) error {
	db := tools.GetDBFromContext(c)
	user := c.Get("db_user").(*models.User)

	var todos []*models.Todo
	db.Model(&models.Todo{}).Where("user_id = ? AND parent_task_id is null", &user.ID).Preload("ChildTasks", models.TodoPreloadChilds).Find(&todos)

	return c.JSON(http.StatusOK, echo.Map{
		"todos": filterTodos(todos),
	})

}

func filterTodos(todos []*models.Todo) []*models.Todo {
	filteredTodos := make([]*models.Todo, 0)

	for _, todo := range todos {
		filteredTodo := &models.Todo{
			ID:              todo.ID,
			TaskName:        todo.TaskName,
			TaskDescription: todo.TaskDescription,
			DueDate:         todo.DueDate,
			Completed:       todo.Completed,
			ChildTasks:      filterTodos(todo.ChildTasks),
			CreatedAt:       todo.CreatedAt,
		}

		filteredTodos = append(filteredTodos, filteredTodo)
	}

	return filteredTodos
}
