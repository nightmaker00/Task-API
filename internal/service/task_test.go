package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/nightmaker00/go-tasks-api/internal/domain"
	mockservice "github.com/nightmaker00/go-tasks-api/internal/service/mock"
	"go.uber.org/mock/gomock"
)

func TestTaskServiceCreate_InvalidTitle_ReturnsErrInvalidTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()

	_, err := service.Create(ctx, "   ", "desc")
	if !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestTaskServiceCreate_TrimsTitleAndNormalizesDescription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()

	repo.EXPECT().
		Create(ctx, gomock.Any(), "My task", nil, string(domain.TaskStatusNew)).
		Return(nil)

	id, err := service.Create(ctx, "  My task  ", "   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == uuid.Nil {
		t.Fatalf("expected non-nil id")
	}
}

func TestTaskServiceGetByID_NotFound_ReturnsErrTaskNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()
	id := uuid.New()

	repo.EXPECT().GetByID(ctx, id).Return(nil, nil)

	_, err := service.GetByID(ctx, id)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskServiceUpdate_InvalidTitle_ReturnsErrInvalidTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()

	err := service.Update(ctx, uuid.New(), "   ", nil, string(domain.TaskStatusNew))
	if !errors.Is(err, ErrInvalidTitle) {
		t.Fatalf("expected ErrInvalidTitle, got %v", err)
	}
}

func TestTaskServiceUpdate_InvalidStatus_ReturnsErrInvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()

	err := service.Update(ctx, uuid.New(), "Title", nil, "bad")
	if !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestTaskServiceUpdate_NotFound_ReturnsErrTaskNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()
	id := uuid.New()

	repo.EXPECT().
		Update(ctx, id, "Title", nil, string(domain.TaskStatusDone)).
		Return(false, nil)

	err := service.Update(ctx, id, "Title", nil, string(domain.TaskStatusDone))
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskServiceUpdate_Success_ReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()
	id := uuid.New()
	desc := "desc"

	repo.EXPECT().
		Update(ctx, id, "Title", &desc, string(domain.TaskStatusNew)).
		Return(true, nil)

	err := service.Update(ctx, id, "Title", &desc, string(domain.TaskStatusNew))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskServiceDelete_Success_ReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()
	id := uuid.New()

	repo.EXPECT().Delete(ctx, id).Return(nil)

	err := service.Delete(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskServiceList_InvalidStatus_ReturnsErrInvalidStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()

	_, err := service.List(ctx, "bad", 10, 0)
	if !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestTaskServiceList_DefaultLimit_Applied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()

	repo.EXPECT().
		List(ctx, "", defaultListLimit, 5).
		Return([]domain.TaskListItem{}, nil)

	_, err := service.List(ctx, "", 0, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskServiceList_InvalidOffset_ReturnsErrInvalidOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mockservice.NewMockTaskRepository(ctrl)
	service := NewTaskService(repo)
	ctx := context.Background()

	_, err := service.List(ctx, "", 10, -1)
	if !errors.Is(err, ErrInvalidOffset) {
		t.Fatalf("expected ErrInvalidOffset, got %v", err)
	}
}

