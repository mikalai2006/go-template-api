package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
)

func (h *Handler) RegisterComponent(router *gin.RouterGroup) {
		route := router.Group("/component")
		{
			route.GET("/:id", h.getComponent)
			route.GET("/find", h.findComponent)
			route.POST("/", middleware.SetUserIdentity, h.createComponent)
			route.DELETE("/:id", middleware.SetUserIdentity, h.deleteComponent)
			route.PATCH("/:id", middleware.SetUserIdentity, h.updateComponent)
		}
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
// @Router /api/component/{id} [get]
func (h *Handler) getComponent(c *gin.Context) {
	id := c.Param("id")

	user, err := h.services.Component.GetComponent(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}


type InputComponent struct {
	domain.RequestParams
	domain.Component
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
// @Router /api/component/find [get]
func (h *Handler) findComponent(c *gin.Context) {
	params, err := utils.GetParamsFromRequest(c, domain.Component{})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	users, err := h.services.Component.FindComponent(params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) createComponent(c *gin.Context) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var input domain.Component
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := h.services.Component.CreateComponent(userId, input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		// utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
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
// @Router /api/component/{id} [delete]
func (h *Handler) deleteComponent(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))

		return
	}

	user, err := h.services.Component.DeleteComponent(id) // , input
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

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
// @Router /api/component/{id} [put]
func (h *Handler) updateComponent(c *gin.Context)  {

	id := c.Param("id")

	var input domain.Component
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err) // .SetMeta(gin.H{"hello": "World"})

		return
	}

	user, err := h.services.Component.UpdateComponent(id, input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, user)
}