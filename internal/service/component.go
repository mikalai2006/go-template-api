package service

import (
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type ComponentService struct {
	repo repository.Component
	i18n config.I18nConfig
}

func NewComponentService(repo repository.Component, i18n config.I18nConfig) *ComponentService {
	return &ComponentService{repo: repo, i18n: i18n}
}

func (s *ComponentService) GetComponent(id string) (domain.Component, error) {
	return s.repo.GetComponent(id)
}

func (s *ComponentService) FindComponent(params domain.RequestParams) (domain.Response[domain.Component], error) {
	return s.repo.FindComponent(params)
}

func (s *ComponentService) FindByPopulate(params domain.RequestParams) (domain.Response[domain.Component], error) {
	return s.repo.FindByPopulate(params)
}
func (s *ComponentService) CreateComponent(
	userID string, component *domain.ComponentCreate,
) (*domain.Component, error) {
	return s.repo.CreateComponent(userID, component)
}

func (s *ComponentService) DeleteComponent(id string) (domain.Component, error) {
	return s.repo.DeleteComponent(id)
}

func (s *ComponentService) UpdateComponent(id string, data interface{}) (domain.Component, error) {
	return s.repo.UpdateComponent(id, data)
}

func (s *ComponentService) FindLibrarys(params domain.RequestParams) (domain.Response[domain.Library], error) {
	return s.repo.FindLibrarys(params)
}
