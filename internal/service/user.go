package service

import (
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id string) (domain.User, error) {
	return s.repo.GetUser(id)
}

func (s *UserService) FindUser(params domain.RequestParams) (domain.Response[domain.User], error) {
	return s.repo.FindUser(params)
}

func (s *UserService) CreateUser(userID string, user *domain.User) (*domain.User, error) {
	return s.repo.CreateUser(userID, user)
}

func (s *UserService) DeleteUser(id string) (domain.User, error) {
	return s.repo.DeleteUser(id)
}

func (s *UserService) UpdateUser(id string, user *domain.User) (domain.User, error) {
	return s.repo.UpdateUser(id, user)
}

func (s *UserService) Iam(userID string) (domain.User, error) {
	return s.repo.Iam(userID)
}
