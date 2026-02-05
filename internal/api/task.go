package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/nightmaker00/go-tasks-api/internal/domain"
)

type TaskService interface {
	Create(ctx context.Context, title string, description string) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	Update(ctx context.Context, id uuid.UUID, title string, description *string, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, status string, limit, offset int) ([]domain.TaskListItem, error)
}
