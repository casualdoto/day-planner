package repository

import (
	"database/sql"

	"day-planner/internal/models"
	"day-planner/internal/util"
)

type CalendarRepository struct {
	db *sql.DB
}

func NewCalendarRepository(db *sql.DB) *CalendarRepository {
	return &CalendarRepository{db: db}
}

func (r *CalendarRepository) Create(input models.CreateCalendarEventInput) (models.CalendarEvent, error) {
	now := util.NowString()
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO calendar_events (title, description, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, input.Title, input.Description, input.StartDate, input.EndDate, now, now).Scan(&id)
	if err != nil {
		return models.CalendarEvent{}, err
	}
	return r.Get(id)
}

func (r *CalendarRepository) Get(id int64) (models.CalendarEvent, error) {
	events, err := r.scan(`SELECT id, title, description, start_date, end_date, completed, completed_at, created_at, updated_at FROM calendar_events WHERE id = $1`, id)
	if err != nil {
		return models.CalendarEvent{}, err
	}
	if len(events) == 0 {
		return models.CalendarEvent{}, sql.ErrNoRows
	}
	return events[0], nil
}

func (r *CalendarRepository) Update(id int64, input models.CreateCalendarEventInput) (models.CalendarEvent, error) {
	_, err := r.db.Exec(`
		UPDATE calendar_events
		SET title = $1, description = $2, start_date = $3, end_date = $4, updated_at = $5
		WHERE id = $6
	`, input.Title, input.Description, input.StartDate, input.EndDate, util.NowString(), id)
	if err != nil {
		return models.CalendarEvent{}, err
	}
	return r.Get(id)
}

func (r *CalendarRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM calendar_events WHERE id = $1`, id)
	return err
}

func (r *CalendarRepository) ToggleCompleted(id int64) (models.CalendarEvent, error) {
	event, err := r.Get(id)
	if err != nil {
		return models.CalendarEvent{}, err
	}
	completed := !event.Completed
	var completedAt *string
	if completed {
		now := util.NowString()
		completedAt = &now
	}
	_, err = r.db.Exec(`UPDATE calendar_events SET completed = $1, completed_at = $2, updated_at = $3 WHERE id = $4`, completed, completedAt, util.NowString(), id)
	if err != nil {
		return models.CalendarEvent{}, err
	}
	return r.Get(id)
}

func (r *CalendarRepository) List(startDate string, endDate string) ([]models.CalendarEvent, error) {
	return r.scan(`
		SELECT id, title, description, start_date, end_date, completed, completed_at, created_at, updated_at
		FROM calendar_events
		WHERE start_date <= $1 AND COALESCE(end_date, start_date) >= $2
		ORDER BY start_date ASC, created_at ASC
	`, endDate, startDate)
}

func (r *CalendarRepository) scan(query string, args ...any) ([]models.CalendarEvent, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []models.CalendarEvent
	for rows.Next() {
		var event models.CalendarEvent
		if err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartDate, &event.EndDate, &event.Completed, &event.CompletedAt, &event.CreatedAt, &event.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
