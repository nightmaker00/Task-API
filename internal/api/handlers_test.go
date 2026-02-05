package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nightmaker00/go-tasks-api/internal/domain"
)

type fakeTaskService struct {
	t          *testing.T
	createFn   func(ctx context.Context, title string, description string) (uuid.UUID, error)
	getByIDFn  func(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	updateFn   func(ctx context.Context, id uuid.UUID, title string, description *string, status string) error
	deleteFn   func(ctx context.Context, id uuid.UUID) error
	listFn     func(ctx context.Context, status string, limit, offset int) ([]domain.TaskListItem, error)
}

func (f *fakeTaskService) Create(ctx context.Context, title string, description string) (uuid.UUID, error) {
	if f.createFn == nil {
		f.t.Fatalf("unexpected Create call")
		return uuid.Nil, nil
	}
	return f.createFn(ctx, title, description)
}

func (f *fakeTaskService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	if f.getByIDFn == nil {
		f.t.Fatalf("unexpected GetByID call")
		return nil, nil
	}
	return f.getByIDFn(ctx, id)
}

func (f *fakeTaskService) Update(ctx context.Context, id uuid.UUID, title string, description *string, status string) error {
	if f.updateFn == nil {
		f.t.Fatalf("unexpected Update call")
		return nil
	}
	return f.updateFn(ctx, id, title, description, status)
}

func (f *fakeTaskService) Delete(ctx context.Context, id uuid.UUID) error {
	if f.deleteFn == nil {
		f.t.Fatalf("unexpected Delete call")
		return nil
	}
	return f.deleteFn(ctx, id)
}

func (f *fakeTaskService) List(ctx context.Context, status string, limit, offset int) ([]domain.TaskListItem, error) {
	if f.listFn == nil {
		f.t.Fatalf("unexpected List call")
		return nil, nil
	}
	return f.listFn(ctx, status, limit, offset)
}

func decodeError(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	return body["error"]
}

func TestCreateTask_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	svc := &fakeTaskService{t: t}
	handler := NewHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("{"))
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if got := decodeError(t, w); got != "invalid json" {
		t.Fatalf("unexpected error message: %s", got)
	}
}

func TestCreateTask_Success_ReturnsCreatedAndID(t *testing.T) {
	taskID := uuid.New()
	svc := &fakeTaskService{
		t: t,
		createFn: func(ctx context.Context, title string, description string) (uuid.UUID, error) {
			if title != "Title" || description != "Desc" {
				t.Fatalf("unexpected input: %q %q", title, description)
			}
			return taskID, nil
		},
	}
	handler := NewHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"Title","description":"Desc"}`))
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var body domain.CreateTaskResponse
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body.ID != taskID {
		t.Fatalf("expected id %s, got %s", taskID, body.ID)
	}
}

func TestGetTask_InvalidID_ReturnsBadRequest(t *testing.T) {
	svc := &fakeTaskService{t: t}
	handler := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/tasks/bad", nil)
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()

	handler.GetTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if got := decodeError(t, w); got != "invalid id" {
		t.Fatalf("unexpected error message: %s", got)
	}
}

func TestGetTask_Success_ReturnsOKAndTask(t *testing.T) {
	id := uuid.New()
	now := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	task := &domain.Task{
		ID:          id,
		Title:       "Title",
		Description: "Desc",
		Status:      domain.TaskStatusNew,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	svc := &fakeTaskService{
		t: t,
		getByIDFn: func(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
			return task, nil
		},
	}
	handler := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/tasks/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	handler.GetTask(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body domain.Task
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body.ID != id || body.Title != task.Title || body.Status != task.Status {
		t.Fatalf("unexpected task in response")
	}
}

func TestUpdateTask_InvalidID_ReturnsBadRequest(t *testing.T) {
	svc := &fakeTaskService{t: t}
	handler := NewHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/tasks/bad", strings.NewReader(`{"title":"t","status":"new"}`))
	req.SetPathValue("id", "bad")
	w := httptest.NewRecorder()

	handler.UpdateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if got := decodeError(t, w); got != "invalid id" {
		t.Fatalf("unexpected error message: %s", got)
	}
}

func TestUpdateTask_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	svc := &fakeTaskService{t: t}
	handler := NewHandler(svc)
	id := uuid.New()

	req := httptest.NewRequest(http.MethodPut, "/tasks/"+id.String(), strings.NewReader("{"))
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	handler.UpdateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if got := decodeError(t, w); got != "invalid json" {
		t.Fatalf("unexpected error message: %s", got)
	}
}

func TestUpdateTask_Success_ReturnsOK(t *testing.T) {
	svc := &fakeTaskService{
		t: t,
		updateFn: func(ctx context.Context, id uuid.UUID, title string, description *string, status string) error {
			return nil
		},
	}
	handler := NewHandler(svc)
	id := uuid.New()

	req := httptest.NewRequest(http.MethodPut, "/tasks/"+id.String(), strings.NewReader(`{"title":"t","status":"new"}`))
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	handler.UpdateTask(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body domain.UpdateTaskResponse
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body.Status != "updated" {
		t.Fatalf("unexpected status: %s", body.Status)
	}
}

func TestDeleteTask_Success_ReturnsNoContent(t *testing.T) {
	svc := &fakeTaskService{
		t: t,
		deleteFn: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}
	handler := NewHandler(svc)
	id := uuid.New()

	req := httptest.NewRequest(http.MethodDelete, "/tasks/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	w := httptest.NewRecorder()

	handler.DeleteTask(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestListTasks_InvalidLimit_ReturnsBadRequest(t *testing.T) {
	svc := &fakeTaskService{t: t}
	handler := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/tasks?limit=abc", nil)
	w := httptest.NewRecorder()

	handler.ListTasks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if got := decodeError(t, w); got != "invalid limit" {
		t.Fatalf("unexpected error message: %s", got)
	}
}

func TestListTasks_Success_ReturnsOKAndItems(t *testing.T) {
	items := []domain.TaskListItem{
		{ID: uuid.New(), Title: "T1", Status: domain.TaskStatusNew},
	}
	svc := &fakeTaskService{
		t: t,
		listFn: func(ctx context.Context, status string, limit, offset int) ([]domain.TaskListItem, error) {
			if status != "new" || limit != 10 || offset != 5 {
				t.Fatalf("unexpected params: %q %d %d", status, limit, offset)
			}
			return items, nil
		},
	}
	handler := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/tasks?status=new&limit=10&offset=5", nil)
	w := httptest.NewRecorder()

	handler.ListTasks(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body []domain.TaskListItem
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(body) != 1 || body[0].ID != items[0].ID {
		t.Fatalf("unexpected items in response")
	}
}
