package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/nightmaker00/go-tasks-api/internal/domain"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, id uuid.UUID, title string, description *string, status string) error {
	desc := toNullString(description)
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tasks (id, title, description, status) VALUES ($1, $2, $3, $4)`,
		id,
		title,
		desc,
		status,
	)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	task := &domain.Task{}
	var description sql.NullString
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1`,
		id,
	).Scan(&task.ID, &task.Title, &description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get task: %w", err)
	}
	task.Description = fromNullString(description)
	return task, nil
}

func (r *TaskRepository) Update(ctx context.Context, id uuid.UUID, title string, description *string, status string) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("update task begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	desc := toNullString(description)
	result, err := tx.ExecContext(
		ctx,
		`UPDATE tasks SET title = $1, description = $2, status = $3, updated_at = NOW() WHERE id = $4`,
		title,
		desc,
		status,
		id,
	)
	if err != nil {
		return false, fmt.Errorf("update task: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("update task rows: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("update task commit: %w", err)
	}
	return affected > 0, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("delete task begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("delete task commit: %w", err)
	}
	return nil
}

func (r *TaskRepository) List(ctx context.Context, status string, limit, offset int) ([]domain.TaskListItem, error) {
	items := make([]domain.TaskListItem, 0)
	query := `SELECT id, title, status FROM tasks`
	args := []any{}
	if status != "" {
		query += ` WHERE status = $1`
		args = append(args, status)
	}
	query += ` ORDER BY id ASC LIMIT $%d OFFSET $%d`
	limitPos := len(args) + 1
	offsetPos := len(args) + 2
	query = fmt.Sprintf(query, limitPos, offsetPos)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.TaskListItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Status); err != nil {
			return nil, fmt.Errorf("scan task list: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate task list: %w", err)
	}
	return items, nil
}

func toNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *value, Valid: true}
}

func fromNullString(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
