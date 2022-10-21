package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/app"
)

func (h *HandlerV1) RegisterComponent(router *gin.RouterGroup) {
	route := router.Group("/component")
	route.GET("/:id", h.getComponent)
	route.GET("/find", h.findComponent)
	route.GET("/populate", h.findByPopulate)
	route.POST("/", middleware.SetUserIdentity, h.createComponent)
	route.DELETE("/:id", middleware.SetUserIdentity, h.deleteComponent)
	route.PATCH("/:id", middleware.SetUserIdentity, h.updateComponent)

	library := router.Group("/library")
	library.GET("/", h.findLibrary)
}

// @Summary Get component by Id
// @Tags component
// @Description get component info
// @ModuleID component
// @Accept  json
// @Produce  json
// @Param id path string true "component id"
// @Success 200 {object} domain.Component
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/component/{id} [get].
func (h *HandlerV1) getComponent(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	user, err := h.services.Component.GetComponent(id)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

type InputComponent struct {
	domain.RequestParams
	// domain.Component
}

// @Summary Find few components
// @Security ApiKeyAuth
// @Tags component
// @Description Input params for search components
// @ModuleID component
// @Accept  json
// @Produce  json
// @Param input query InputComponent true "params for search components"
// @Success 200 {object} []domain.Component
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/component/find [get].
func (h *HandlerV1) findComponent(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Component{}, &h.i18n)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	users, err := h.services.Component.FindComponent(params)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *HandlerV1) createComponent(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.ComponentCreate
	if er := c.BindJSON(&input); er != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	user, err := h.services.Component.CreateComponent(userID, input)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		// utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Find few components and populate
// @Security ApiKeyAuth
// @Tags component
// @Description Input params for search components
// @ModuleID component
// @Accept  json
// @Produce  json
// @Param input query InputComponent true "params for search components"
// @Success 200 {object} []domain.Component
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/component/populate [get].
func (h *HandlerV1) findByPopulate(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Component{}, &h.i18n)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	users, err := h.services.Component.FindByPopulate(params)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Delete component
// @Security ApiKeyAuth
// @Tags component
// @Description Delete component
// @ModuleID component
// @Accept  json
// @Produce  json
// @Param id path string true "component id"
// @Success 200 {object} domain.Component
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/component/{id} [delete].
func (h *HandlerV1) deleteComponent(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.Component.DeleteComponent(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update component
// @Security ApiKeyAuth
// @Tags component
// @Description Update component
// @ModuleID component
// @Accept  json
// @Produce  json
// @Param id path string true "component id"
// @Param input body domain.Component true "body for update component"
// @Success 200 {object} domain.Component
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/component/{id} [put].
func (h *HandlerV1) updateComponent(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.PageInputData
	data, err := utils.BindAndValid(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	user, err := h.services.Component.UpdateComponent(id, data)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Find few library and populate
// @Security ApiKeyAuth
// @Tags library
// @Description Input params for search librarys
// @ModuleID library
// @Accept  json
// @Produce  json
// @Param input query InputComponent true "params for search librarys"
// @Success 200 {object} []domain.Library
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/library/ [get].
func (h *HandlerV1) findLibrary(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.LibraryInput{}, &h.i18n)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	users, err := h.services.Component.FindLibrarys(params)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}
