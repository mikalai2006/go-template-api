package service

import (
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type ComponentPresetService struct {
	repo repository.ComponentPreset
	i18n config.I18nConfig
}

func NewComponentPresetService(repo repository.ComponentPreset, i18n config.I18nConfig) *ComponentPresetService {
	return &ComponentPresetService{repo: repo, i18n: i18n}
}

func (s *ComponentPresetService) FindComponentPreset(params domain.RequestParams) (domain.Response[domain.ComponentPreset], error) {
	return s.repo.FindComponentPreset(params)
}

func (s *ComponentPresetService) CreateComponentPreset(
	userID string, component_preset *domain.ComponentPresetInput,
) (*domain.ComponentPreset, error) {
	return s.repo.CreateComponentPreset(userID, component_preset)
}

func (s *ComponentPresetService) UpdateComponentPreset(id string, data interface{}) (domain.ComponentPreset, error) {
	return s.repo.UpdateComponentPreset(id, data)
}

func (s *ComponentPresetService) DeleteComponentPreset(id string) (domain.ComponentPreset, error) {
	return s.repo.DeleteComponentPreset(id)
}
