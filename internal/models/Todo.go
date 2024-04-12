package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	ID uint64 `gorm:"primary_key"`

	TaskName        string       `json:"task_name"`
	TaskDescription *string      `json:"task_description,omitempty"`
	DueDate         sql.NullTime `json:"due_date"`
	Completed       bool         `gorm:"default:0" json:"completed"`

	ParentTaskID *uint64 `gorm:"index"`
	ParentTask   *Todo   `gorm:"foreignKey:ParentTaskID"`

	UserId *uint64 `json:"index"`
	User   *User   `gorm:"foreignKey:UserId"`
}
