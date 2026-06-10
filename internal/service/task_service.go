package service

import (
	"errors"
	"strings"

	"day-planner/internal/models"
	"day-planner/internal/repository"
	"day-planner/internal/util"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(input models.CreateTaskInput) (models.Task, error) {
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		return models.Task{}, errors.New("task title is required")
	}
	if err := util.ValidateDate(input.Date); err != nil {
		return models.Task{}, err
	}
	if input.DueDate != nil && *input.DueDate != "" {
		if err := util.ValidateDate(*input.DueDate); err != nil {
			return models.Task{}, err
		}
	}
	if input.Importance == 0 {
		input.Importance = 2
	}
	if input.Importance < 1 || input.Importance > 4 {
		return models.Task{}, errors.New("importance must be between 1 and 4")
	}
	return s.repo.Create(input)
}

func (s *TaskService) UpdateTask(id int64, input models.UpdateTaskInput) (models.Task, error) {
	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return models.Task{}, errors.New("task title is required")
		}
		input.Title = &title
	}
	if input.Date != nil {
		if err := util.ValidateDate(*input.Date); err != nil {
			return models.Task{}, err
		}
	}
	if input.DueDate != nil && *input.DueDate != "" {
		if err := util.ValidateDate(*input.DueDate); err != nil {
			return models.Task{}, err
		}
	}
	if input.Importance != nil && (*input.Importance < 1 || *input.Importance > 4) {
		return models.Task{}, errors.New("importance must be between 1 and 4")
	}
	return s.repo.Update(id, input)
}

func (s *TaskService) DeleteTask(id int64) error {
	return s.repo.Delete(id)
}

func (s *TaskService) ToggleTaskCompleted(id int64) (models.Task, error) {
	return s.repo.ToggleCompleted(id)
}

func (s *TaskService) ListTasks(filter models.TaskFilter) ([]models.Task, error) {
	if filter.Date != nil && *filter.Date != "" {
		if err := util.ValidateDate(*filter.Date); err != nil {
			return nil, err
		}
	}
	return s.repo.List(filter)
}
