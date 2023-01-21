package service

import (
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type PartnerService struct {
	repo repository.Partner
	i18n config.I18nConfig
}

func NewPartnerService(repo repository.Partner, i18n config.I18nConfig) *PartnerService {
	return &PartnerService{repo: repo, i18n: i18n}
}

func (s *PartnerService) CreatePartner(userID string, data *domain.PartnerInput) (domain.Partner, error) {
	return s.repo.CreatePartner(userID, data)
}

func (s *PartnerService) GetPartner(id string) (domain.Partner, error) {
	return s.repo.GetPartner(id)
}

func (s *PartnerService) FindPartner(params domain.RequestParams) (domain.Response[domain.Partner], error) {
	return s.repo.FindPartner(params)
}

func (s *PartnerService) UpdatePartner(id string, data interface{}) (domain.Partner, error) {
	return s.repo.UpdatePartner(id, data)
}
func (s *PartnerService) DeletePartner(id string) (domain.Partner, error) {
	return s.repo.DeletePartner(id)
}
