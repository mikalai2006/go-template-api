package service

import (
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type PageService struct {
	repo repository.Page
}

func NewPageService(repo repository.Page) *PageService  {
	return &PageService{repo: repo}
}

func (s *PageService) GetPage(id string) (domain.Page, error) {
	return s.repo.GetPage(id)
}

func (s *PageService) FindPage(params domain.RequestParams) (domain.Response[domain.Page], error) {
	return s.repo.FindPage(params)
}

func (s *PageService) CreatePage(userId string, page domain.Page) (*domain.Page, error)  {
	return s.repo.CreatePage(userId, page)
}

func (s *PageService) DeletePage(id string) (domain.Page, error) {
	return s.repo.DeletePage(id)
}

func (s *PageService) UpdatePage(id string, page domain.Page) (domain.Page, error) {
	return s.repo.UpdatePage(id, page)
}
