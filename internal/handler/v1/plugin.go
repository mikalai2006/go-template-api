package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/app"
)

// func init() {
// 	if _, err := os.Stat("public/css"); errors.Is(err, os.ErrNotExist) {
// 		err := os.MkdirAll("public/css", os.ModePerm)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

func (h HandlerV1) RegisterPlugin(router *gin.RouterGroup) {
	route := router.Group("/plugin")
	route.POST("", middleware.SetUserIdentity, h.createPlugin)
	route.GET("/:id", h.getPlugin)
	route.GET("", h.findPlugin)
	route.PATCH("/:id", middleware.SetUserIdentity, h.updatePlugin)
	route.DELETE("/:id", middleware.SetUserIdentity, h.deletePlugin)
}

func (h HandlerV1) createPlugin(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.PluginInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	Plugin, err := h.services.Plugin.CreatePlugin(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Plugin)
}

func (h HandlerV1) getPlugin(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	document, err := h.services.Plugin.GetPlugin(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h HandlerV1) findPlugin(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.PluginInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	results, err := h.services.Plugin.FindPlugin(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h HandlerV1) updatePlugin(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.PluginInput
	data, err := utils.BindAndValid(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	fmt.Println("update Plugin", data)

	result, err := h.services.Plugin.UpdatePlugin(id, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h HandlerV1) deletePlugin(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	document, err := h.services.Plugin.DeletePlugin(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}
