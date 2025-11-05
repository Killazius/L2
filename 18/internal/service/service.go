package service

import (
	"18/internal/domain"
	"errors"
	"time"
)

type Repository interface {
	Create(e domain.Event) (domain.Event, error)
	Update(e domain.Event) error
	Delete(userID string, date string, id int64) error
	GetByDate(userID string, date string) ([]domain.Event, error)
	GetByWeek(userID string, date string) ([]domain.Event, error)
	GetByMonth(userID string, date string) ([]domain.Event, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

var ErrInvalidDate = errors.New("invalid date format, expected YYYY-MM-DD")

func validateDate(date string) error {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return ErrInvalidDate
	}
	return nil
}
func (s *Service) Create(e domain.Event) (domain.Event, error) {
	if err := validateDate(e.Date); err != nil {
		return domain.Event{}, err
	}
	return s.repo.Create(e)
}

func (s *Service) Update(e domain.Event) error {
	if err := validateDate(e.Date); err != nil {
		return err
	}
	return s.repo.Update(e)
}

func (s *Service) Delete(userID string, date string, id int64) error {
	if err := validateDate(date); err != nil {
		return err
	}

	return s.repo.Delete(userID, date, id)
}

func (s *Service) GetByDate(userID string, date string) ([]domain.Event, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}
	return s.repo.GetByDate(userID, date)
}

func (s *Service) GetByWeek(userID string, date string) ([]domain.Event, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}
	return s.repo.GetByWeek(userID, date)
}

func (s *Service) GetByMonth(userID string, date string) ([]domain.Event, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}
	return s.repo.GetByMonth(userID, date)
}
