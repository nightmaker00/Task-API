package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nightmaker00/go-tasks-api/internal/domain"
)

//go:generate mockgen -package=service -source=repository.go -destination=./mock/mock.go
type TaskRepository interface {
	Create(ctx context.Context, id uuid.UUID, title string, description *string, status string) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	Update(ctx context.Context, id uuid.UUID, title string, description *string, status string) (bool, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, status string, limit, offset int) ([]domain.TaskListItem, error)
}
