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

func (h *HandlerV1) RegisterComponentPreset(router *gin.RouterGroup) {
	route := router.Group("/component_preset")
	route.GET("/", h.findComponentPreset)
	// route.GET("/:id", h.findComponentPreset)
	route.POST("/", middleware.SetUserIdentity, h.createComponentPreset)
	route.PATCH("/:id", middleware.SetUserIdentity, h.updateComponentPreset)
	route.DELETE("/:id", middleware.SetUserIdentity, h.deleteComponentPreset)
}

func (h *HandlerV1) findComponentPreset(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.ComponentPresetFind{}, &h.i18n)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	presets, err := h.services.ComponentPreset.FindComponentPreset(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, presets)
}

func (h *HandlerV1) createComponentPreset(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.ComponentPresetInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	user, err := h.services.ComponentPreset.CreateComponentPreset(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) updateComponentPreset(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.ComponentPresetInput
	err := c.BindJSON(&input) // utils.BindAndValid(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	user, err := h.services.ComponentPreset.UpdateComponentPreset(id, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) deleteComponentPreset(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.ComponentPreset.DeleteComponentPreset(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
