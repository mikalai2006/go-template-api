package service

import (
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type AppsService struct {
	repo repository.Apps
	i18n config.I18nConfig
}

func NewAppsService(repo repository.Apps, i18n config.I18nConfig) *AppsService {
	return &AppsService{repo: repo, i18n: i18n}
}

func (s *AppsService) CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error) {
	return s.repo.CreateLanguage(userID, data)
}

func (s *AppsService) GetLanguage(id string) (domain.Language, error) {
	return s.repo.GetLanguage(id)
}

func (s *AppsService) FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error) {
	return s.repo.FindLanguage(params)
}

func (s *AppsService) UpdateLanguage(id string, data interface{}) (domain.Language, error) {
	return s.repo.UpdateLanguage(id, data)
}
func (s *AppsService) DeleteLanguage(id string) (domain.Language, error) {
	return s.repo.DeleteLanguage(id)
}
