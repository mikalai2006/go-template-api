package service

import (
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type ComponentService struct {
	repo repository.Component
}

func NewComponentService(repo repository.Component) *ComponentService  {
	return &ComponentService{repo: repo}
}

func (s *ComponentService) GetComponent(id string) (domain.Component, error) {
	return s.repo.GetComponent(id)
}

func (s *ComponentService) FindComponent(params domain.RequestParams) (domain.Response[domain.Component], error) {
	return s.repo.FindComponent(params)
}

func (s *ComponentService) CreateComponent(userId string, component domain.Component) (*domain.Component, error)  {
	return s.repo.CreateComponent(userId, component)
}

func (s *ComponentService) DeleteComponent(id string) (domain.Component, error) {
	return s.repo.DeleteComponent(id)
}

func (s *ComponentService) UpdateComponent(id string, component domain.Component) (domain.Component, error) {
	return s.repo.UpdateComponent(id, component)
}
