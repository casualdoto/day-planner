package db

import "database/sql"

func Migrate(conn *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS day_entries (
			id BIGSERIAL PRIMARY KEY,
			date TEXT NOT NULL UNIQUE,
			plan_text TEXT NOT NULL DEFAULT '',
			result_text TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tags (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			color TEXT NOT NULL DEFAULT '#64748b'
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			date TEXT NOT NULL,
			due_date TEXT,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			completed_at TEXT,
			importance INTEGER NOT NULL DEFAULT 2,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS task_tags (
			task_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			PRIMARY KEY (task_id, tag_id),
			FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS calendar_events (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			start_date TEXT NOT NULL,
			end_date TEXT,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			completed_at TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_date ON tasks(date)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(completed)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_importance ON tasks(importance)`,
		`CREATE INDEX IF NOT EXISTS idx_calendar_events_start_date ON calendar_events(start_date)`,
	}
	for _, statement := range statements {
		if _, err := conn.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}
