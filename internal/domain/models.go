package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusNew        TaskStatus = "new"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

// Task представляет задачу
// @Description Задача с UUID, заголовком, описанием, статусом и временными метками
type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TaskListItem представляет краткую информацию о задаче в списке
// @Description Краткая информация о задаче для списка
type TaskListItem struct {
	ID     uuid.UUID  `json:"id"`
	Title  string     `json:"title"`
	Status TaskStatus `json:"status"`
}

// CreateTaskRequest запрос на создание задачи
// @Description Данные для создания новой задачи
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTaskRequest запрос на обновление задачи
// @Description Данные для обновления задачи
type UpdateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
}

// CreateTaskResponse ответ при создании задачи
// @Description UUID созданной задачи
type CreateTaskResponse struct {
	ID uuid.UUID `json:"id"`
}

// UpdateTaskResponse ответ при обновлении задачи
// @Description Статус операции обновления
type UpdateTaskResponse struct {
	Status string `json:"status"`
}
