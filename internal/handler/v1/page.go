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

func (h HandlerV1) RegisterPage(router *gin.RouterGroup) {
	page := router.Group("/page")
	page.POST("/", middleware.SetUserIdentity, h.createPage)
	page.DELETE("/:id", middleware.SetUserIdentity, h.deletePage)
	page.PATCH("/:id", middleware.SetUserIdentity, h.updatePage)
	page.GET("/:id", h.getPage)
	page.GET("/get", h.getFullPage)
	page.GET("/find", h.findPage)
	page.GET("/routers", h.getPageForRouters)
}

// @Summary Get pages by routers
// @Tags page
// @Description Get pages by routers
// @ModuleID page
// @Accept  json
// @Produce  json
// @Success 200 {object} domain.PageRoutes
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/page [get].
func (h HandlerV1) getPageForRouters(c *gin.Context) {
	appG := app.Gin{C: c}

	document, err := h.services.Page.GetPageForRouters()
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

// @Summary Get page by slug
// @Tags page
// @Description get page info
// @ModuleID page
// @Accept  json
// @Produce  json
// @Param slug path string true "page slug"
// @Success 200 {object} domain.Page
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/page/{slug} [get].
func (h HandlerV1) getFullPage(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.PageInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	document, err := h.services.Page.GetFullPage(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

// @Summary Get page by Id
// @Tags page
// @Description get page info
// @ModuleID page
// @Accept  json
// @Produce  json
// @Param id path string true "page id"
// @Success 200 {object} domain.Page
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/page/{id} [get].
func (h HandlerV1) getPage(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	document, err := h.services.Page.GetPage(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

type InputPage struct {
	domain.RequestParams
	domain.Page
}

// @Summary Find few pages
// @Security ApiKeyAuth
// @Tags page
// @Description Input params for search pages
// @ModuleID page
// @Accept  json
// @Produce  json
// @Param input query InputPage true "params for search pages"
// @Success 200 {object} []domain.Page
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/page/find [get].
func (h HandlerV1) findPage(c *gin.Context) {
	appG := app.Gin{C: c}

	// var params domain.PageQuery
	// if err := c.Bind(&params); err != nil {
	// 	appG.Response(http.StatusBadRequest, err, nil)
	// 	return
	// }

	params, err := utils.GetParamsFromRequest(c, domain.PageInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	documents, err := h.services.Page.FindPage(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (h HandlerV1) createPage(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *domain.PageInputData
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.services.Page.CreatePage(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

// @Summary Delete page
// @Security ApiKeyAuth
// @Tags page
// @Description Delete page
// @ModuleID page
// @Accept  json
// @Produce  json
// @Param id path string true "page id"
// @Success 200 {object} domain.Page
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/page/{id} [delete].
func (h HandlerV1) deletePage(c *gin.Context) {
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

	document, err := h.services.Page.DeletePage(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

// @Summary Update page
// @Security ApiKeyAuth
// @Tags page
// @Description Update page
// @ModuleID page
// @Accept  json
// @Produce  json
// @Param id path string true "page id"
// @Param input body domain.Page true "body for update page"
// @Success 200 {object} domain.Page
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/page/{id} [put].
func (h HandlerV1) updatePage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.PageInputData
	data, err := utils.BindAndValid(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	document, err := h.services.Page.UpdatePage(id, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}
