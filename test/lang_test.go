package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/repository"
)

var testLanuageData = domain.Language{
	Name:      "testLang",
	Code:      "ru",
	Flag:      "flag",
	Publish:   true,
	Locale:    "ru-Ru",
	SortOrder: 1,
}

func (s *TestSuite) TestCreateLangNotAuth() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	dataJSON, err := json.Marshal(testLanuageData)

	// test invalid header
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/lang/", bytes.NewBuffer(dataJSON))
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer")
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	// test empty token
	req, err = http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/lang/", bytes.NewBuffer(dataJSON))
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer ")
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()

	r.Equal(http.StatusUnauthorized, response.StatusCode)
}

func (s *TestSuite) TestCreateLangAuth() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))

	auth, err := s.Auth(router)
	s.NoError(err)

	coll := s.db.Collection(repository.TblLanguage)
	err = coll.Drop(context.Background())
	s.NoError(err)

	r := s.Require()

	dataJSON, err := json.Marshal(testLanuageData)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/lang/", bytes.NewBuffer(dataJSON))
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var re domain.Language
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.Code, testLanuageData.Code)
	r.Equal(re.Flag, testLanuageData.Flag)
	r.Equal(re.Name, testLanuageData.Name)
	r.Equal(re.Publish, testLanuageData.Publish)
	r.Equal(re.SortOrder, testLanuageData.SortOrder)
	r.Equal(re.Locale, testLanuageData.Locale)

	r.Equal(http.StatusOK, response.StatusCode)
}

func (s *TestSuite) TestFindLangByLimitOne() {
	limit := 1

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/lang/", nil)
	q := req.URL.Query()
	q.Add("$limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	req.Close = true
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var re domain.Response[domain.Language]
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.Limit, limit)
	r.Equal(len(re.Data), limit)

	r.Equal(http.StatusOK, response.StatusCode)
}
func (s *TestSuite) TestFindLangByLimitBig() {
	limit := 1000

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/lang/", nil)
	q := req.URL.Query()
	q.Add("$limit", fmt.Sprintf("%d", limit))
	req.URL.RawQuery = q.Encode()
	req.Close = true
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var re domain.Response[domain.Language]
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.Limit, 10)
	r.Equal(len(re.Data), re.Total)

	r.Equal(http.StatusOK, response.StatusCode)
}

func (s *TestSuite) TestFindLangBySkip() {
	skip := 1

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/lang/", nil)
	q := req.URL.Query()
	q.Add("$skip", fmt.Sprintf("%d", skip))
	req.URL.RawQuery = q.Encode()
	req.Close = true
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var re domain.Response[domain.Language]
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.Skip, skip)
	r.Equal(len(re.Data), re.Total)

	r.Equal(http.StatusOK, response.StatusCode)
}

func (s *TestSuite) TestFindLangBySort() {
	sort := -1

	router := gin.New()
	s.handler.Init(router.Group("/api"))

	// create two item.
	auth, err := s.Auth(router)
	s.NoError(err)

	testLanuageDataTwo := domain.Language{
		Name:      "Two",
		SortOrder: 10,
		Publish:   true,
	}
	dataJSON, err := json.Marshal(testLanuageDataTwo)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/lang/", bytes.NewBuffer(dataJSON))
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var twoItem domain.Language
	err = json.NewDecoder(response.Body).Decode(&twoItem)
	s.NoError(err)

	r := s.Require()

	req, err = http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/lang/", nil)
	q := url.Values{}
	q.Add("$sort[sort_order]", fmt.Sprintf("%v", sort))
	req.URL.RawQuery = q.Encode()
	req.Close = true
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()

	var re domain.Response[domain.Language]
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	r.Equal(re.Data[0].SortOrder, testLanuageDataTwo.SortOrder)
	// r.Equal(len(re.Data), re.Total)

	r.Equal(http.StatusOK, response.StatusCode)
}

func (s *TestSuite) TestDeleteLang() {
	router := gin.New()
	s.handler.Init(router.Group("/api"))

	// create item.
	auth, err := s.Auth(router)
	s.NoError(err)

	testData := domain.Language{
		Name:      "For remove",
		SortOrder: 1,
		Publish:   false,
	}
	dataJSON, err := json.Marshal(testData)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/lang/", bytes.NewBuffer(dataJSON))
	req.Close = true
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	s.NoError(err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response := w.Result()
	defer response.Body.Close()

	var twoItem domain.Language
	err = json.NewDecoder(response.Body).Decode(&twoItem)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)

	r := s.Require()

	// test empty id.
	req, err = http.NewRequestWithContext(context.Background(), http.MethodDelete, "/api/v1/lang/ ", nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Close = true
	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()
	r.Equal(http.StatusBadRequest, response.StatusCode)

	// test with id.
	idForRemove := twoItem.ID.Hex()
	req, err = http.NewRequestWithContext(context.Background(), http.MethodDelete, fmt.Sprintf("/api/v1/lang/%s", idForRemove), nil)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth.AccessToken)
	req.Close = true

	s.NoError(err)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	response = w.Result()
	defer response.Body.Close()

	var re domain.Language
	err = json.NewDecoder(response.Body).Decode(&re)
	s.NoError(err)

	//r.Equal(re.ID.Hex(), twoItem.ID.Hex())
	// r.Equal(len(re.Data), re.Total)

	r.Equal(http.StatusOK, response.StatusCode)
}
