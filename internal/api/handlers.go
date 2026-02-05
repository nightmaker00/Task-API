package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/nightmaker00/go-tasks-api/internal/domain"
	"github.com/nightmaker00/go-tasks-api/internal/service"
)

type Handler struct {
	taskService TaskService
}

func NewHandler(taskService TaskService) *Handler {
	return &Handler{taskService: taskService}
}

// CreateTask создаёт новую задачу
// @Summary      Создать задачу
// @Description  Создаёт новую задачу с указанным заголовком и описанием
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body      domain.CreateTaskRequest  true  "Данные задачи"
// @Success      201   {object}  domain.CreateTaskResponse
// @Failure      400   {object}  map[string]string  "Неверный запрос"
// @Failure      500   {object}  map[string]string  "Внутренняя ошибка"
// @Router       /tasks [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	id, err := h.taskService.Create(r.Context(), req.Title, req.Description)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, domain.CreateTaskResponse{ID: id})
}

// GetTask получает задачу по ID
// @Summary      Получить задачу
// @Description  Возвращает задачу по её UUID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "UUID задачи"
// @Success      200  {object}  domain.Task
// @Failure      400  {object}  map[string]string  "Неверный UUID"
// @Failure      404  {object}  map[string]string  "Задача не найдена"
// @Router       /tasks/{id} [get]
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	task, err := h.taskService.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toTaskResponse(task))
}

// UpdateTask обновляет задачу
// @Summary      Обновить задачу
// @Description  Обновляет данные задачи (заголовок, описание, статус)
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    path      string                true  "UUID задачи"
// @Param        task  body      domain.UpdateTaskRequest  true  "Обновлённые данные"
// @Success      200   {object}  domain.UpdateTaskResponse
// @Failure      400   {object}  map[string]string  "Неверный запрос"
// @Failure      404   {object}  map[string]string  "Задача не найдена"
// @Router       /tasks/{id} [put]
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req domain.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	err = h.taskService.Update(r.Context(), id, req.Title, req.Description, req.Status)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, domain.UpdateTaskResponse{Status: "updated"})
}

// DeleteTask удаляет задачу
// @Summary      Удалить задачу
// @Description  Удаляет задачу по её UUID (идемпотентная операция)
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "UUID задачи"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string  "Неверный UUID"
// @Router       /tasks/{id} [delete]
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.taskService.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

// ListTasks получает список задач
// @Summary      Список задач
// @Description  Возвращает список задач с фильтрацией по статусу и пагинацией
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        status  query     string  false  "Фильтр по статусу (new, in_progress, done)"
// @Param        limit   query     int     false  "Лимит записей (по умолчанию 100, максимум 1000)"
// @Param        offset  query     int     false  "Смещение для пагинации"
// @Success      200     {array}   domain.TaskListItem
// @Failure      400     {object}  map[string]string  "Неверные параметры"
// @Router       /tasks [get]
func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	limit, err := parseIntParam(r, "limit")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid limit")
		return
	}
	offset, err := parseIntParam(r, "offset")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid offset")
		return
	}

	items, err := h.taskService.List(r.Context(), status, limit, offset)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toTaskListResponse(items))
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrTaskNotFound):
		writeError(w, http.StatusNotFound, "task not found")
	case errors.Is(err, service.ErrInvalidTitle),
		errors.Is(err, service.ErrInvalidStatus),
		errors.Is(err, service.ErrInvalidLimit),
		errors.Is(err, service.ErrInvalidOffset):
		writeError(w, http.StatusBadRequest, "invalid request")
	default:
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", h.GetTask)
	mux.HandleFunc("PUT /tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)
	mux.HandleFunc("GET /tasks", h.ListTasks)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func parseID(raw string) (uuid.UUID, error) {
	value, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return uuid.Nil, errors.New("invalid uuid")
	}
	return value, nil
}

func parseIntParam(r *http.Request, key string) (int, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func toTaskResponse(task *domain.Task) *domain.Task {
	if task == nil {
		return nil
	}
	return task
}

func toTaskListResponse(items []domain.TaskListItem) []domain.TaskListItem {
	return items
}
