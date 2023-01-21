package service

import (
	"errors"
	"fmt"
	"os"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type PluginService struct {
	repo repository.Plugin
	i18n config.I18nConfig
}

func NewPluginService(repo repository.Plugin, i18n config.I18nConfig) *PluginService {
	return &PluginService{repo: repo, i18n: i18n}
}

func (s *PluginService) CreatePlugin(userID string, plugin *domain.PluginInput) (*domain.Plugin, error) {

	result, err := s.repo.CreatePlugin(userID, plugin)
	if err != nil {
		return result, err
	}

	fmt.Println("plugin data:", plugin.Code, plugin.SpaceID)
	if plugin.Code != "" && plugin.SpaceID != "" {
		dirPath := fmt.Sprintf("./public/s_%s/pl", plugin.SpaceID)

		fmt.Println("plugin path:", dirPath)

		if _, err := os.Stat(dirPath); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				return result, err
			}
		}

		filePath := fmt.Sprintf("%s/%s.js", dirPath, result.ID.Hex())
		f, err := os.Create(filePath)
		if err != nil {
			return result, err
		}
		defer f.Close()

		f.Write([]byte(plugin.Code))
	}

	return result, err
}

func (s *PluginService) GetPlugin(id string) (domain.Plugin, error) {
	return s.repo.GetPlugin(id)
}

func (s *PluginService) FindPlugin(params domain.RequestParams) (domain.Response[domain.Plugin], error) {
	return s.repo.FindPlugin(params)
}

func (s *PluginService) UpdatePlugin(id string, data interface{}) (domain.Plugin, error) {
	return s.repo.UpdatePlugin(id, data)
}

func (s *PluginService) DeletePlugin(id string) (domain.Plugin, error) {
	return s.repo.DeletePlugin(id)
}
