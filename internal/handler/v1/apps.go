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

func (h *HandlerV1) RegisterApp(router *gin.RouterGroup) {
	lang := router.Group("/lang")
	lang.GET("/", h.findLanguage)
	lang.GET("/:id", h.getLanguage)
	lang.POST("/", middleware.SetUserIdentity, h.createLanguage)
	lang.PATCH("/:id", middleware.SetUserIdentity, h.updateLanguage)
	lang.DELETE("/:id", middleware.SetUserIdentity, h.deleteLanguage)
}

func (h *HandlerV1) createLanguage(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		appG.Response(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.LanguageInput
	if er := c.BindJSON(&input); er != nil {
		appG.Response(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.services.Apps.CreateLanguage(userID, input)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) getLanguage(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	document, err := h.services.Apps.GetLanguage(id)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) findLanguage(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.LanguageInput{})
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	documents, err := h.services.Apps.FindLanguage(params)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (h *HandlerV1) updateLanguage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.LanguageInput
	data, err := app.BindAndValid(c, &input)
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}
	fmt.Println(data)

	document, err := h.services.Apps.UpdateLanguage(id, &data)
	if err != nil {
		appG.Response(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) deleteLanguage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.Response(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}
	// var input domain.Page
	// if err := c.BindJSON(&input); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)

	// 	return
	// }

	document, err := h.services.Apps.DeleteLanguage(id) // , input
	if err != nil {
		appG.Response(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}
