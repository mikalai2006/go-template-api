package service

import (
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type ComponentGroupService struct {
	repo repository.ComponentGroup
	i18n config.I18nConfig
}

func NewComponentGroupService(repo repository.ComponentGroup, i18n config.I18nConfig) *ComponentGroupService {
	return &ComponentGroupService{repo: repo, i18n: i18n}
}

func (s *ComponentGroupService) FindComponentGroup() (domain.Response[domain.ComponentGroup], error) {
	return s.repo.FindComponentGroup()
}

func (s *ComponentGroupService) CreateComponentGroup(
	userID string, component_group *domain.ComponentGroup,
) (*domain.ComponentGroup, error) {
	return s.repo.CreateComponentGroup(userID, component_group)
}

func (s *ComponentGroupService) UpdateComponentGroup(id string, data interface{}) (domain.ComponentGroup, error) {
	return s.repo.UpdateComponentGroup(id, data)
}

func (s *ComponentGroupService) DeleteComponentGroup(id string) (domain.ComponentGroup, error) {
	return s.repo.DeleteComponentGroup(id)
}
