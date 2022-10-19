package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/app"
)

func (h *HandlerV1) RegisterUser(router *gin.RouterGroup) {
	user := router.Group("/user")
	user.POST("/", middleware.SetUserIdentity, h.CreateUser)
	user.DELETE("/:id", middleware.SetUserIdentity, h.DeleteUser)
	user.PATCH("/:id", middleware.SetUserIdentity, h.UpdateUser)
	user.GET("/:id", h.GetUser)
	user.GET("/find", h.FindUser)
}

// @Summary Get user by Id
// @Tags user
// @Description get user info
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} domain.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [get].
func (h *HandlerV1) GetUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	user, err := h.services.User.GetUser(id)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// type InputUser struct {
// 	domain.RequestParams
// 	domain.User
// }

// @Summary Find few users
// @Security ApiKeyAuth
// @Tags user
// @Description Input params for search users
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param input query domain.UserInput true "params for search users"
// @Success 200 {object} []domain.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user [get].
func (h *HandlerV1) FindUser(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.User{})
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	users, err := h.services.User.FindUser(params)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *HandlerV1) CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		appG.Response(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.User
	if er := c.BindJSON(&input); er != nil {
		appG.Response(http.StatusBadRequest, er, nil)
		return
	}

	user, er := h.services.User.CreateUser(userID, input)
	if er != nil {
		appG.Response(http.StatusBadRequest, er, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Security ApiKeyAuth
// @Tags user
// @Description Delete user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} []domain.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [delete].
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	// var input domain.User
	// if err := c.BindJSON(&input); err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())

	// 	return
	// }

	user, err := h.services.User.DeleteUser(id) // , input
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update user
// @Security ApiKeyAuth
// @Tags user
// @Description Update user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Param input body domain.User true "body for update user"
// @Success 200 {object} []domain.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [put].
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input *domain.User
	if err := c.ShouldBindJSON(&input); err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	user, err := h.services.User.UpdateUser(id, input)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
