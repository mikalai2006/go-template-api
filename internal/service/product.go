package service

import (
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type ProductService struct {
	repo repository.Product
	i18n config.I18nConfig
}

func NewProductService(repo repository.Product, i18n config.I18nConfig) *ProductService {
	return &ProductService{repo: repo, i18n: i18n}
}

func (s *ProductService) CreateProduct(userID string, data *domain.ProductInput) (domain.Product, error) {
	return s.repo.CreateProduct(userID, data)
}

func (s *ProductService) GetProduct(id string) (domain.Product, error) {
	return s.repo.GetProduct(id)
}

func (s *ProductService) FindProduct(params domain.RequestParams) (domain.Response[domain.Product], error) {
	return s.repo.FindProduct(params)
}

func (s *ProductService) UpdateProduct(id string, data interface{}) (domain.Product, error) {
	return s.repo.UpdateProduct(id, data)
}
func (s *ProductService) DeleteProduct(id string) (domain.Product, error) {
	return s.repo.DeleteProduct(id)
}
