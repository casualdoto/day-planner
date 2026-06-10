package service

import (
	"day-planner/internal/models"
	"day-planner/internal/repository"
	"day-planner/internal/util"
)

type DayService struct {
	repo *repository.DayRepository
}

func NewDayService(repo *repository.DayRepository) *DayService {
	return &DayService{repo: repo}
}

func (s *DayService) GetDayEntry(date string) (models.DayEntry, error) {
	if err := util.ValidateDate(date); err != nil {
		return models.DayEntry{}, err
	}
	return s.repo.Get(date)
}

func (s *DayService) SaveDayPlan(date string, planText string) error {
	if err := util.ValidateDate(date); err != nil {
		return err
	}
	return s.repo.SavePlan(date, planText)
}

func (s *DayService) SaveDayResult(date string, resultText string) error {
	if err := util.ValidateDate(date); err != nil {
		return err
	}
	return s.repo.SaveResult(date, resultText)
}
