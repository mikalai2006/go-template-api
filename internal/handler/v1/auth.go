package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/pkg/app"
)

func (h *HandlerV1) registerAuth(router *gin.RouterGroup) {
	router.POST("/sign-up", h.SignUp)
	router.POST("/sign-in", h.SignIn)
	router.POST("/logout", h.Logout)
	router.POST("/refresh", h.tokenRefresh)
	router.GET("/refresh", h.tokenRefresh)
	router.GET("/verification/:code", middleware.SetUserIdentity, h.VerificationAuth)
	router.GET("/iam", middleware.SetUserIdentity, h.getIam)
}

func (h *HandlerV1) getIam(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		appG.Response(http.StatusUnauthorized, err, nil)
		return
	}

	// TODO get token from body data.
	// var input *domain.RefreshInput

	// if err := c.BindJSON(&input); err != nil {
	// 	appG.Response(http.StatusBadRequest, err, nil)
	// 	return
	// }

	users, err := h.services.User.Iam(userID)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	// implementation roles for user.
	if roles, err := middleware.GetRoles(c); err != nil {
		appG.Response(http.StatusUnauthorized, err, nil)
		return
	} else {
		users.Roles = roles
	}

	c.JSON(http.StatusOK, users)
}

// @Summary SignUp
// @Tags auth
// @Description Create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body domain.Auth true "account info"
// @Success 200 {integer} 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/sign-up [post].
func (h *HandlerV1) SignUp(c *gin.Context) {
	appG := app.Gin{C: c}

	var input *domain.SignInInput

	if err := c.BindJSON(&input); err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	id, err := h.services.Authorization.CreateAuth(input)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary SignIn
// @Tags auth
// @Description Login user
// @ID signin-account
// @Accept json
// @Produce json
// @Param input body domain.SignInInput true "credentials"
// @Success 200 {integer} 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/sign-in [post].
func (h *HandlerV1) SignIn(c *gin.Context) {
	appG := app.Gin{C: c}
	// jwt_cookie, _ := c.Cookie("jwt-handmade")
	// fmt.Println("+++++++++++++")
	// fmt.Printf("jwt_handmade = %s", jwt_cookie)
	// fmt.Println("+++++++++++++")
	// session := sessions.Default(c)
	var input *domain.SignInInput

	if err := c.BindJSON(&input); err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	if input.Strategy == "" {
		input.Strategy = "local"
	}

	if input.Email == "" && input.Login == "" {
		appG.Response(http.StatusBadRequest, errors.New("request must be with email or login"), nil)
		return
	}

	if input.Strategy == "local" {
		tokens, err := h.services.Authorization.SignIn(input)
		if err != nil {
			appG.Response(http.StatusBadRequest, err, nil)
			return
		}
		c.SetCookie("jwt-handmade", tokens.RefreshToken, h.oauth.TimeExpireCookie, "/", c.Request.URL.Hostname(), false, true)

		c.JSON(http.StatusOK, domain.ResponseTokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
	// else {
	// 	fmt.Print("JWT auth")
	// }
	// session.Set(userkey, input.Username)
	// session.Save()
}

// @Summary User Refresh Tokens
// @Tags users-auth
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body domain.RefreshInput true "sign up info"
// @Success 200 {object} domain.ResponseTokens
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /users/auth/refresh [post].
func (h *HandlerV1) tokenRefresh(c *gin.Context) {
	appG := app.Gin{C: c}
	jwtCookie, _ := c.Cookie("jwt-handmade")
	// fmt.Println("jwt_handmade = ", jwtCookie)
	// jwt_header := c.GetHeader("hello")
	// fmt.Println("jwt_header = ", jwt_header)
	// fmt.Println("+++++++++++++")
	// session := sessions.Default(c)
	var input domain.RefreshInput

	if jwtCookie == "" {
		if err := c.BindJSON(&input); err != nil {
			appG.Response(http.StatusBadRequest, err, nil)
			return
		}
	} else {
		input.Token = jwtCookie
	}

	if input.Token == "" && jwtCookie == "" {
		c.JSON(http.StatusOK, gin.H{})
		c.AbortWithStatus(http.StatusOK)
		return
	}

	res, err := h.services.Authorization.RefreshTokens(input.Token)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.SetCookie("jwt-handmade", res.RefreshToken, h.oauth.TimeExpireCookie, "/", c.Request.URL.Hostname(), false, true)

	c.JSON(http.StatusOK, domain.ResponseTokens{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	})
}

func (h *HandlerV1) Logout(c *gin.Context) {
	// session := sessions.Default(c)
	// session.Delete(userkey)
	// if err := session.Save(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func (h *HandlerV1) VerificationAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	code := c.Param("code")
	if code == "" {
		appG.Response(http.StatusBadRequest, errors.New("code empty"), nil)
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	if er := h.services.Authorization.VerificationCode(userID, code); er != nil {
		appG.Response(http.StatusBadRequest, er, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}