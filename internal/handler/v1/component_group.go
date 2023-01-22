package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/pkg/app"
)

func (h *HandlerV1) RegisterComponentGroup(router *gin.RouterGroup) {
	route := router.Group("/component_group")
	route.GET("/", h.findComponentGroup)
	route.POST("/", middleware.SetUserIdentity, h.createComponentGroup)
	route.PATCH("/:id", middleware.SetUserIdentity, h.updateComponentGroup)
	route.DELETE("/:id", middleware.SetUserIdentity, h.deleteComponentGroup)
}

func (h *HandlerV1) findComponentGroup(c *gin.Context) {
	appG := app.Gin{C: c}

	groups, err := h.services.ComponentGroup.FindComponentGroup()
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (h *HandlerV1) createComponentGroup(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.ComponentGroup
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	user, err := h.services.ComponentGroup.CreateComponentGroup(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) updateComponentGroup(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.ComponentGroupInput
	err := c.BindJSON(&input) // utils.BindAndValid(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	user, err := h.services.ComponentGroup.UpdateComponentGroup(id, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) deleteComponentGroup(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.services.ComponentGroup.DeleteComponentGroup(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
