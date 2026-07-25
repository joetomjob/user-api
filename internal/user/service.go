package user

import (
	"context"
	"errors"
	"strings"
)

type UserRepository interface {
	Create(ctx context.Context, user User) (User, error)
	GetById(ctx context.Context, id int) (User, error)
	Update(ctx context.Context, user User, id int) (User, error)
	Delete(ctx context.Context, id int) error
}

type Service struct {
	repo UserRepository
}

var ErrNameRequired = errors.New("User name is required")
var ErrEmailRequired = errors.New("Email is required")
var ErrInvalidAge = errors.New("Invalid Age")
var ErrInvalidId = errors.New("Invalid Id")

func NewService(repo UserRepository) *Service {
	return &Service{repo}
}

func (s *Service) Create(ctx context.Context, user User) (User, error) {
	if strings.TrimSpace(user.Name) == "" {
		return User{}, ErrNameRequired
	}

	if strings.TrimSpace(user.Email) == "" {
		return User{}, ErrEmailRequired
	}

	if user.Age < 0 {
		return User{}, ErrInvalidAge
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
		return User{}, ErrInvalidId
	}

	if strings.TrimSpace(user.Name) == "" {
		return User{}, ErrNameRequired
	}

	if strings.TrimSpace(user.Email) == "" {
		return User{}, ErrEmailRequired
	}

	if user.Age < 0 {
		return User{}, ErrInvalidAge
	}

	user, err := s.repo.Update(ctx, user, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id < 0 {
		return ErrInvalidId
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
