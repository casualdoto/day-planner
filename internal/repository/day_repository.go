package repository

import (
	"database/sql"

	"day-planner/internal/models"
	"day-planner/internal/util"
)

type DayRepository struct {
	db *sql.DB
}

func NewDayRepository(db *sql.DB) *DayRepository {
	return &DayRepository{db: db}
}

func (r *DayRepository) Get(date string) (models.DayEntry, error) {
	row := r.db.QueryRow(`SELECT id, date, plan_text, result_text, created_at, updated_at FROM day_entries WHERE date = $1`, date)
	var entry models.DayEntry
	err := row.Scan(&entry.ID, &entry.Date, &entry.PlanText, &entry.ResultText, &entry.CreatedAt, &entry.UpdatedAt)
	if err == sql.ErrNoRows {
		now := util.NowString()
		return models.DayEntry{Date: date, CreatedAt: now, UpdatedAt: now}, nil
	}
	return entry, err
}

func (r *DayRepository) SavePlan(date string, text string) error {
	now := util.NowString()
	_, err := r.db.Exec(`
		INSERT INTO day_entries (date, plan_text, result_text, created_at, updated_at)
		VALUES ($1, $2, '', $3, $4)
		ON CONFLICT(date) DO UPDATE SET plan_text = excluded.plan_text, updated_at = excluded.updated_at
	`, date, text, now, now)
	return err
}

func (r *DayRepository) SaveResult(date string, text string) error {
	now := util.NowString()
	_, err := r.db.Exec(`
		INSERT INTO day_entries (date, plan_text, result_text, created_at, updated_at)
		VALUES ($1, '', $2, $3, $4)
		ON CONFLICT(date) DO UPDATE SET result_text = excluded.result_text, updated_at = excluded.updated_at
	`, date, text, now, now)
	return err
}
