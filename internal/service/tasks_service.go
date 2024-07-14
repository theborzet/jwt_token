package service

import (
	"errors"
	"time"

	"github.com/theborzet/time-tracker/internal/models"
	"github.com/theborzet/time-tracker/internal/repository"
)

func (s *ApiService) GetUserTasks(userId int, start, end string) ([]*models.TaskTimeSpent, error) {
	if start == "" || end == "" {
		return nil, errors.New("start or end date cannot be empty")
	}
	if userId <= 0 {
		return nil, errors.New("incorrect userId value")
	}

	startTime, err := time.Parse(repository.DateFormat, start)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse(repository.DateFormat, end)
	if err != nil {
		return nil, err
	}

	startTimeFormatted := startTime.Format(repository.TimestampFormat)
	endTimeFormatted := endTime.Format(repository.TimestampFormat)
	tasks, err := s.repo.GetUserTasks(userId, startTimeFormatted, endTimeFormatted)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *ApiService) StartTask(userId int, taskName string) error {
	if taskName == "" {
		return errors.New("taskName cannot be empty")
	}
	if userId <= 0 {
		return errors.New("incorrect userId value")
	}

	startTime := time.Now()
	if err := s.repo.StartTask(userId, taskName, startTime.Format(repository.TimestampFormat)); err != nil {
		return err
	}
	return nil
}

func (s *ApiService) EndTask(userId int, taskName string) error {
	if taskName == "" {
		return errors.New("taskName cannot be empty")
	}
	if userId <= 0 {
		return errors.New("incorrect userId value")
	}
	endTime := time.Now()
	if err := s.repo.StartTask(userId, taskName, endTime.Format(repository.TimestampFormat)); err != nil {
		return err
	}
	return nil
}
