package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/pkg/app"
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func (h *HandlerV1) registerGoogleOAuth(router *gin.RouterGroup) {
	router.GET("/google", h.OAuthGoogle)
	router.GET("/google/me", h.MeGoogle)
}

func (h *HandlerV1) OAuthGoogle(c *gin.Context) {
	appG := app.Gin{C: c}

	urlReferer := c.Request.Referer()
	scope := strings.Join(h.oauth.GoogleScopes, " ")

	pathRequest, err := url.Parse(h.oauth.GoogleAuthURI)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	parameters := url.Values{}
	parameters.Add("client_id", h.oauth.GoogleClientID)
	parameters.Add("redirect_uri", h.oauth.GoogleRedirectURI)
	parameters.Add("scope", scope)
	parameters.Add("response_type", "code")
	parameters.Add("state", urlReferer)

	pathRequest.RawQuery = parameters.Encode()
	c.Redirect(http.StatusFound, pathRequest.String())
}

func (h *HandlerV1) MeGoogle(c *gin.Context) {
	appG := app.Gin{C: c}

	code := c.Query("code")
	clientURL := c.Query("state")

	if code == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("no correct code"))
		appG.Response(http.StatusBadRequest, errors.New("no correct code"), nil)
		return
	}

	pathRequest, err := url.Parse(h.oauth.GoogleTokenURI)
	if err != nil {
		panic("boom")
	}
	parameters := url.Values{}
	parameters.Set("client_id", h.oauth.GoogleClientID)
	parameters.Set("redirect_uri", h.oauth.GoogleRedirectURI)
	parameters.Set("client_secret", h.oauth.GoogleClientSecret)
	parameters.Set("code", code)
	parameters.Set("grant_type", "authorization_code")

	req, _ := http.NewRequestWithContext(c, http.MethodPost, pathRequest.String(), strings.NewReader(parameters.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}
	defer resp.Body.Close()

	token := struct {
		AccessToken string `json:"access_token"`
	}{}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}
	if er := json.Unmarshal(bytes, &token); er != nil { // Parse []byte to go struct pointer
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, er, nil)
		return
	}

	pathRequest, err = url.Parse(h.oauth.GoogleUserinfoURI)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}
	r, _ := http.NewRequestWithContext(c, http.MethodGet, pathRequest.String(), http.NoBody) // URL-encoded payload
	bearerToken := fmt.Sprintf("Bearer %s", token.AccessToken)
	r.Header.Add("Authorization", bearerToken)
	// r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}
	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	var bodyResponse GoogleUserInfo
	if e := json.Unmarshal(bytes, &bodyResponse); e != nil { // Parse []byte to go struct pointer
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, e, nil)
		return
	}

	input := &domain.SignInInput{
		Login:    bodyResponse.Email,
		Strategy: "jwt",
		Password: "",
		GoogleID: bodyResponse.Sub,
	}

	user, err := h.services.Authorization.ExistAuth(input)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	if user.Login == "" {
		_, err = h.services.Authorization.CreateAuth(input)
		if err != nil {
			// c.AbortWithError(http.StatusBadRequest, err)
			appG.Response(http.StatusBadRequest, err, nil)
			return
		}
	}

	tokens, err := h.services.Authorization.SignIn(input)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	pathRequest, err = url.Parse(clientURL)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}
	parameters = url.Values{}
	parameters.Add("token", tokens.AccessToken)
	pathRequest.RawQuery = parameters.Encode()
	c.SetCookie("jwt-handmade", tokens.RefreshToken, h.oauth.TimeExpireCookie, "/", c.Request.URL.Hostname(), false, true)
	c.Redirect(http.StatusFound, pathRequest.String())
}
