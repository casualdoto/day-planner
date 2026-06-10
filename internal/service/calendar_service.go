package service

import (
	"errors"
	"strings"

	"day-planner/internal/models"
	"day-planner/internal/repository"
	"day-planner/internal/util"
)

type CalendarService struct {
	repo *repository.CalendarRepository
}

func NewCalendarService(repo *repository.CalendarRepository) *CalendarService {
	return &CalendarService{repo: repo}
}

func (s *CalendarService) CreateCalendarEvent(input models.CreateCalendarEventInput) (models.CalendarEvent, error) {
	if err := validateEvent(input); err != nil {
		return models.CalendarEvent{}, err
	}
	return s.repo.Create(input)
}

func (s *CalendarService) UpdateCalendarEvent(id int64, input models.CreateCalendarEventInput) (models.CalendarEvent, error) {
	if err := validateEvent(input); err != nil {
		return models.CalendarEvent{}, err
	}
	return s.repo.Update(id, input)
}

func (s *CalendarService) DeleteCalendarEvent(id int64) error {
	return s.repo.Delete(id)
}

func (s *CalendarService) ToggleCalendarEventCompleted(id int64) (models.CalendarEvent, error) {
	return s.repo.ToggleCompleted(id)
}

func (s *CalendarService) ListCalendarEvents(startDate string, endDate string) ([]models.CalendarEvent, error) {
	if err := util.ValidateDate(startDate); err != nil {
		return nil, err
	}
	if err := util.ValidateDate(endDate); err != nil {
		return nil, err
	}
	return s.repo.List(startDate, endDate)
}

func validateEvent(input models.CreateCalendarEventInput) error {
	if strings.TrimSpace(input.Title) == "" {
		return errors.New("event title is required")
	}
	if err := util.ValidateDate(input.StartDate); err != nil {
		return err
	}
	if input.EndDate != nil && *input.EndDate != "" {
		return util.ValidateDate(*input.EndDate)
	}
	return nil
}
