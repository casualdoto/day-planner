package repository

import (
	"database/sql"

	"day-planner/internal/models"
)

type TagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(name string, color string) (models.Tag, error) {
	var id int64
	if err := r.db.QueryRow(`INSERT INTO tags (name, color) VALUES ($1, $2) RETURNING id`, name, color).Scan(&id); err != nil {
		return models.Tag{}, err
	}
	return r.Get(id)
}

func (r *TagRepository) Get(id int64) (models.Tag, error) {
	var tag models.Tag
	err := r.db.QueryRow(`SELECT id, name, color FROM tags WHERE id = $1`, id).Scan(&tag.ID, &tag.Name, &tag.Color)
	return tag, err
}

func (r *TagRepository) Update(id int64, name string, color string) (models.Tag, error) {
	_, err := r.db.Exec(`UPDATE tags SET name = $1, color = $2 WHERE id = $3`, name, color, id)
	if err != nil {
		return models.Tag{}, err
	}
	return r.Get(id)
}

func (r *TagRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tags WHERE id = $1`, id)
	return err
}

func (r *TagRepository) List() ([]models.Tag, error) {
	rows, err := r.db.Query(`SELECT id, name, color FROM tags ORDER BY LOWER(name)`)
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
