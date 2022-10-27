package service

import (
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type PageService struct {
	repo repository.Page
	i18n config.I18nConfig
}

func NewPageService(repo repository.Page, i18n config.I18nConfig) *PageService {
	return &PageService{repo: repo, i18n: i18n}
}

func (s *PageService) GetPageForRouters() (domain.Response[domain.PageRoutes], error) {
	return s.repo.GetPageForRouters()
}

func (s *PageService) GetPage(id string) (domain.Page, error) {
	return s.repo.GetPage(id)
}

func (s *PageService) GetFullPage(params domain.RequestParams) (domain.Response[domain.Page], error) {
	return s.repo.GetFullPage(params)
}

func (s *PageService) FindPage(params domain.RequestParams) (domain.Response[domain.Page], error) {
	return s.repo.FindPage(params)
}

func (s *PageService) CreatePage(userID string, page *domain.PageInputData) (*domain.Page, error) {
	return s.repo.CreatePage(userID, page)
}

func (s *PageService) DeletePage(id string) (domain.Page, error) {
	return s.repo.DeletePage(id)
}

func (s *PageService) UpdatePage(id string, data interface{}) (domain.Page, error) {
	return s.repo.UpdatePage(id, data)
}
