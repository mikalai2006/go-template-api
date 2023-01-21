package service

import (
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type SpaceService struct {
	repo repository.Space
	i18n config.I18nConfig
}

func NewSpaceService(repo repository.Space, i18n config.I18nConfig) *SpaceService {
	return &SpaceService{repo: repo, i18n: i18n}
}

func (s *SpaceService) CreateSpace(userID string, space *domain.SpaceInput) (*domain.Space, error) {
	return s.repo.CreateSpace(userID, space)
}

func (s *SpaceService) GetSpace(id string) (domain.Space, error) {
	return s.repo.GetSpace(id)
}

func (s *SpaceService) FindSpace(params domain.RequestParams) (domain.Response[domain.Space], error) {
	return s.repo.FindSpace(params)
}

func (s *SpaceService) UpdateSpace(id string, data interface{}) (domain.Space, error) {
	return s.repo.UpdateSpace(id, data)
}

func (s *SpaceService) DeleteSpace(id string) (domain.Space, error) {
	return s.repo.DeleteSpace(id)
}

// func (s *PageService) GetPageForRouters() (domain.Response[domain.PageRoutes], error) {
// 	return s.repo.GetPageForRouters()
// }

// func (s *PageService) GetFullPage(params domain.RequestParams) (domain.Response[domain.Page], error) {
// 	return s.repo.GetFullPage(params)
// }

// func (s *PageService) UpdatePageWithContent(id string, data map[string]interface{}) (domain.Page, error) {
// 	return s.repo.UpdatePageWithContent(id, data)
// }
