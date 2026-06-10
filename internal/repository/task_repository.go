package repository

import (
	"database/sql"
	"strconv"
	"strings"

	"day-planner/internal/models"
	"day-planner/internal/util"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(input models.CreateTaskInput) (models.Task, error) {
	now := util.NowString()
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO tasks (title, description, date, due_date, importance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, input.Title, input.Description, input.Date, input.DueDate, input.Importance, now, now).Scan(&id)
	if err != nil {
		return models.Task{}, err
	}
	if err := r.replaceTags(id, input.TagIDs); err != nil {
		return models.Task{}, err
	}
	return r.Get(id)
}

func (r *TaskRepository) Get(id int64) (models.Task, error) {
	tasks, err := r.scanTasks(`SELECT id, title, description, date, due_date, completed, completed_at, importance, created_at, updated_at FROM tasks WHERE id = $1`, id)
	if err != nil {
		return models.Task{}, err
	}
	if len(tasks) == 0 {
		return models.Task{}, sql.ErrNoRows
	}
	return tasks[0], nil
}

func (r *TaskRepository) Update(id int64, input models.UpdateTaskInput) (models.Task, error) {
	current, err := r.Get(id)
	if err != nil {
		return models.Task{}, err
	}
	title := current.Title
	description := current.Description
	date := current.Date
	dueDate := current.DueDate
	completed := current.Completed
	completedAt := current.CompletedAt
	importance := current.Importance
	if input.Title != nil {
		title = *input.Title
	}
	if input.Description != nil {
		description = *input.Description
	}
	if input.Date != nil {
		date = *input.Date
	}
	if input.DueDate != nil {
		dueDate = input.DueDate
	}
	if input.Completed != nil {
		completed = *input.Completed
		if completed && completedAt == nil {
			now := util.NowString()
			completedAt = &now
		}
		if !completed {
			completedAt = nil
		}
	}
	if input.Importance != nil {
		importance = *input.Importance
	}

	_, err = r.db.Exec(`
		UPDATE tasks
		SET title = $1, description = $2, date = $3, due_date = $4, completed = $5, completed_at = $6, importance = $7, updated_at = $8
		WHERE id = $9
	`, title, description, date, dueDate, completed, completedAt, importance, util.NowString(), id)
	if err != nil {
		return models.Task{}, err
	}
	if input.TagIDs != nil {
		if err := r.replaceTags(id, input.TagIDs); err != nil {
			return models.Task{}, err
		}
	}
	return r.Get(id)
}

func (r *TaskRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	return err
}

func (r *TaskRepository) ToggleCompleted(id int64) (models.Task, error) {
	task, err := r.Get(id)
	if err != nil {
		return models.Task{}, err
	}
	next := !task.Completed
	return r.Update(id, models.UpdateTaskInput{Completed: &next})
}

func (r *TaskRepository) List(filter models.TaskFilter) ([]models.Task, error) {
	query := `SELECT DISTINCT t.id, t.title, t.description, t.date, t.due_date, t.completed, t.completed_at, t.importance, t.created_at, t.updated_at FROM tasks t`
	args := []any{}
	wheres := []string{}
	nextPlaceholder := func(value any) string {
		args = append(args, value)
		return "$" + strconv.Itoa(len(args))
	}
	if filter.TagID != nil {
		query += ` JOIN task_tags tt ON tt.task_id = t.id`
		wheres = append(wheres, `tt.tag_id = `+nextPlaceholder(*filter.TagID))
	}
	if filter.Date != nil && *filter.Date != "" {
		wheres = append(wheres, `t.date = `+nextPlaceholder(*filter.Date))
	}
	if filter.Importance != nil {
		wheres = append(wheres, `t.importance = `+nextPlaceholder(*filter.Importance))
	}
	if filter.Completed != nil {
		wheres = append(wheres, `t.completed = `+nextPlaceholder(*filter.Completed))
	}
	if filter.Query != nil && strings.TrimSpace(*filter.Query) != "" {
		needle := "%" + strings.ToLower(strings.TrimSpace(*filter.Query)) + "%"
		titlePlaceholder := nextPlaceholder(needle)
		descriptionPlaceholder := nextPlaceholder(needle)
		wheres = append(wheres, `(LOWER(t.title) LIKE `+titlePlaceholder+` OR LOWER(t.description) LIKE `+descriptionPlaceholder+`)`)
	}
	if len(wheres) > 0 {
		query += ` WHERE ` + strings.Join(wheres, ` AND `)
	}
	query += ` ORDER BY t.completed ASC, t.importance DESC, t.created_at DESC`
	return r.scanTasks(query, args...)
}

func (r *TaskRepository) scanTasks(query string, args ...any) ([]models.Task, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Date, &task.DueDate, &task.Completed, &task.CompletedAt, &task.Importance, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for i := range tasks {
		tags, err := r.tagsForTask(tasks[i].ID)
		if err != nil {
			return nil, err
		}
		tasks[i].Tags = tags
	}
	return tasks, nil
}

func (r *TaskRepository) replaceTags(taskID int64, tagIDs []int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM task_tags WHERE task_id = $1`, taskID); err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, tagID := range tagIDs {
		if _, err := tx.Exec(`INSERT INTO task_tags (task_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, taskID, tagID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *TaskRepository) tagsForTask(taskID int64) ([]models.Tag, error) {
	rows, err := r.db.Query(`
		SELECT tags.id, tags.name, tags.color
		FROM tags
		JOIN task_tags ON task_tags.tag_id = tags.id
		WHERE task_tags.task_id = $1
		ORDER BY LOWER(tags.name)
	`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.Color); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}
