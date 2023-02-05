package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

type StoryService struct {
	repo repository.Story
	i18n config.I18nConfig
}

func NewStoryService(repo repository.Story, i18n config.I18nConfig) *StoryService {
	return &StoryService{repo: repo, i18n: i18n}
}

// func (s *StoryService) GetStory(params domain.RequestParams) (domain.Story, error) {
// 	return s.repo.GetStory(params)
// }

// func (s *StoryService) FindStory(params domain.RequestParams) (domain.Response[domain.Story], error) {
// 	return s.repo.FindStory(params)
// }

// func (s *StoryService) CreateStory(userID string, Story *domain.StoryInputData) (*domain.Story, error) {
// 	return s.repo.CreateStory(userID, Story)
// }

// func (s *StoryService) DeleteStory(id string) (domain.Story, error) {
// 	return s.repo.DeleteStory(id)
// }

func (s *StoryService) PublishStory(id string, data domain.StoryInputData) (domain.Story, error) {

	dataStory := buildStoryContent(data.Content)

	data.Content = dataStory

	fmt.Println("data.Content")
	// if val, ok := data.Content; ok {
	// 	fmt.Println("val", val)
	// 	// data := val.([]domain.ComponentData)
	// 	// var ui []interface{}
	// 	// for _, t := range data {
	// 	// 	ui = append(ui, t)
	// 	// }
	// 	// _, err := collection.InsertMany(ctx, ui)
	// 	// if err != nil {
	// 	// 	return result, err
	// 	// }

	// }

	return s.repo.PublishStory(id, data)
}

func buildStoryContent(
	content interface{},
	// i18n config.I18nConfig,
) map[string]interface{} {
	result := map[string]interface{}{}

	cont := content.(map[string]interface{})
	for key, _ := range cont {
		// fmt.Println("+++++++++++++++++key=", key, reflect.TypeOf(cont[key]))
		if reflect.ValueOf(cont[key]).Kind() == reflect.Slice {
			// fmt.Println("slice=", key, reflect.TypeOf(cont[key]))

			nested := []interface{}{}

			arr := cont[key].([]interface{})
			for keySlice, _ := range arr {
				// fmt.Println("Kind slice=", reflect.ValueOf(arr[keySlice]).Kind())

				if reflect.ValueOf(arr[keySlice]).Kind() == reflect.Map {
					nested = append(nested, buildStoryContent(arr[keySlice]))
				} else {
					nested = append(nested, arr[keySlice])
				}
			}
			result[key] = nested
		} else if key == "_uid" {
			b := make([]byte, 4) //equals 8 characters
			rand.Read(b)
			s := hex.EncodeToString(b)

			result["_uid"] = s
		} else {
			if (key != "parent" && key != "global" && key != "_uid") || (key == "parent" && cont[key] == "page") {
				// fmt.Println("value=", key)
				result[key] = cont[key]
			}
		}

	}

	// fmt.Println("============result======================")
	// fmt.Println(result)
	// fmt.Println("==========================================")

	return result
}

// func (s *StoryService) UpdateStoryWithContent(id string, data map[string]interface{}) (domain.Story, error) {
// 	return s.repo.UpdateStoryWithContent(id, data)
// }

func (s *StoryService) GetStory(params domain.RequestParams) (domain.Story, error) {
	return s.repo.GetStory(params)
}
