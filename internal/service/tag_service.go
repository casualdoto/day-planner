package service

import (
	"errors"
	"strings"

	"day-planner/internal/models"
	"day-planner/internal/repository"
)

type TagService struct {
	repo *repository.TagRepository
}

func NewTagService(repo *repository.TagRepository) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) CreateTag(name string, color string) (models.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Tag{}, errors.New("tag name is required")
	}
	if strings.TrimSpace(color) == "" {
		color = "#64748b"
	}
	return s.repo.Create(name, color)
}

func (s *TagService) UpdateTag(id int64, name string, color string) (models.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Tag{}, errors.New("tag name is required")
	}
	if strings.TrimSpace(color) == "" {
		color = "#64748b"
	}
	return s.repo.Update(id, name, color)
}

func (s *TagService) DeleteTag(id int64) error {
	return s.repo.Delete(id)
}

func (s *TagService) ListTags() ([]models.Tag, error) {
	return s.repo.List()
}
