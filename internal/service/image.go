package service

import (
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type ImageService struct {
	repo repository.Image
}

func NewImageService(repo repository.Image) *ImageService {
	return &ImageService{repo: repo}
}

func (s *ImageService) FindImage(params domain.RequestParams) (domain.Response[domain.Image], error) {
	return s.repo.FindImage(params)
}

func (s *ImageService) GetImage(id string) (domain.Image, error) {
	return s.repo.GetImage(id)
}

func (s *ImageService) CreateImage(userID string, image *domain.ImageInput) (domain.Image, error) {
	return s.repo.CreateImage(userID, image)
}

func (s *ImageService) DeleteImage(id string) (domain.Image, error) {
	return s.repo.DeleteImage(id)
}
