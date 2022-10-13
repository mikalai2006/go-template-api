package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
)

func (h *Handler) RegisterPage(router *gin.RouterGroup) {
		page := router.Group("/page")
		{
			page.POST("/", middleware.SetUserIdentity, h.createPage)
			page.DELETE("/:id", middleware.SetUserIdentity, h.deletePage)
			page.PATCH("/:id", middleware.SetUserIdentity, h.updatePage)
			page.GET("/:id", h.getPage)
			page.GET("/find", h.findPage)
		}
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
// @Router /api/page/{id} [get]
func (h *Handler) getPage(c *gin.Context) {
	id := c.Param("id")

	document, err := h.services.Page.GetPage(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
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
// @Router /api/page/find [get]
func (h *Handler) findPage(c *gin.Context) {
	params, err := utils.GetParamsFromRequest(c, domain.Page{})
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	documents, err := h.services.Page.FindPage(params)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (h *Handler) createPage(c *gin.Context) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var input domain.Page
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	document, err := h.services.Page.CreatePage(userId, input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		// utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
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
// @Router /api/page/{id} [delete]
func (h *Handler) deletePage(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))

		return
	}
	// var input domain.Page
	// if err := c.BindJSON(&input); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, err)

	// 	return
	// }

	document, err := h.services.Page.DeletePage(id) // , input
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

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
// @Router /api/page/{id} [put]
func (h *Handler) updatePage(c *gin.Context)  {

	id := c.Param("id")

	var input domain.Page
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithError(http.StatusBadRequest, err) // .SetMeta(gin.H{"hello": "World"})

		return
	}

	document, err := h.services.Page.UpdatePage(id, input)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, document)
}