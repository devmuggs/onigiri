package users

import (
	"context"
	"fmt"
)

type Service interface {
	GetUser(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, user *CreateInput) error
	UpdateUser(ctx context.Context, id int64, user *User) error
	GetAllUsers(ctx context.Context) ([]*User, error)
	GetUserByEmail(ctx context.Context, email string) (*UserRecord, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) Service {
	return &userService{repo: repo}
}

func (s *userService) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, user *CreateInput) error {
	if user.Email == "" {
		return fmt.Errorf("username required")
	}

	return s.repo.Create(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, id int64, user *User) error {
	if user.Email == "" {
		return fmt.Errorf("username required")
	}

	return s.repo.Update(ctx, id, user)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*UserRecord, error) {
	return s.repo.FindUserByEmail(ctx, email)
}
