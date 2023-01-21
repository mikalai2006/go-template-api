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

func (h HandlerV1) RegisterSpace(router *gin.RouterGroup) {
	page := router.Group("/space")
	page.POST("", middleware.SetUserIdentity, h.createSpace)
	page.GET("/:id", h.getSpace)
	page.GET("", h.findSpace)
	page.PATCH("/:id", middleware.SetUserIdentity, h.updateSpace)
	page.DELETE("/:id", middleware.SetUserIdentity, h.deleteSpace)
	// page.PATCH("/:id/content", middleware.SetUserIdentity, h.updatePageWithContent)
	// page.GET("/get", h.getFullPage)
}

func (h HandlerV1) createSpace(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.SpaceInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	space, err := h.services.Space.CreateSpace(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, space)
}

func (h HandlerV1) getSpace(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	document, err := h.services.Space.GetSpace(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h HandlerV1) findSpace(c *gin.Context) {
	appG := app.Gin{C: c}

	// var params domain.PageQuery
	// if err := c.Bind(&params); err != nil {
	// 	appG.Response(http.StatusBadRequest, err, nil)
	// 	return
	// }

	params, err := utils.GetParamsFromRequest(c, domain.SpaceInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	results, err := h.services.Space.FindSpace(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h HandlerV1) updateSpace(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.SpaceInput
	data, err := utils.BindAndValid(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	fmt.Println("update space", data)

	result, err := h.services.Space.UpdateSpace(id, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h HandlerV1) deleteSpace(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}
	// var input domain.Page
	// if err := c.BindJSON(&input); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)

	// 	return
	// }

	document, err := h.services.Space.DeleteSpace(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}
