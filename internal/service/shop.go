package service

import (
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type ShopService struct {
	repo repository.Shop
}

func NewShopService(repo repository.Shop) *ShopService {
	return &ShopService{repo: repo}
}

func (s *ShopService) FindShop(params domain.RequestParams) (domain.Response[domain.Shop], error) {
	return s.repo.FindShop(params)
}

func (s *ShopService) GetAllShops(params domain.RequestParams) (domain.Response[domain.Shop], error) {
	return s.repo.GetAllShops(params)
}

func (s *ShopService) CreateShop(userID string, shop *domain.Shop) (*domain.Shop, error) {
	return s.repo.CreateShop(userID, shop)
}
