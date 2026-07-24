package user

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo}
}

func (s *Service) Create(ctx context.Context, user User) (User, error) {
	if strings.TrimSpace(user.Name) == "" {
		return User{}, errors.New("User name is required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return User{}, errors.New("Email is required")
	}

	if user.Age < 0 {
		return User{}, errors.New("Age should be greater than or equal to zero")
	}

	user, err := s.repo.Create(ctx, user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Service) GetById(ctx context.Context, id int) (User, error) {
	if id < 0 {
		return User{}, errors.New("invalid id")
	}

	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Service) Update(ctx context.Context, user User, id int) (User, error) {
	if id < 0 {
		return User{}, errors.New("Invalid Id")
	}

	if strings.TrimSpace(user.Name) == "" {
		return User{}, errors.New("User name is required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return User{}, errors.New("Email is required")
	}

	if user.Age < 0 {
		return User{}, errors.New("Age should be greater than or equal to zero")
	}

	user, err := s.repo.Update(ctx, user, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id < 0 {
		return errors.New("invalid id")
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
