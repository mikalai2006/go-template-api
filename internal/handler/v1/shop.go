package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/app"
)

func (h *HandlerV1) registerShop(router *gin.RouterGroup) {
	shops := router.Group("/shops")
	shops.GET("/", h.FindShop)
	shops.POST("/", middleware.SetUserIdentity, h.CreateShop)
}

func (h *HandlerV1) CreateShop(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *domain.Shop
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	shop, err := h.services.Shop.CreateShop(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, shop)
}

// @Summary Shop Get all shops
// @Security ApiKeyAuth
// @Tags shop
// @Description get all shops
// @ModuleID shops
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Shop
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/shops [get].
func (h *HandlerV1) GetAllShops(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Shop{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	shops, err := h.services.Shop.GetAllShops(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, shops)
}

// @Summary Find shops by params
// @Security ApiKeyAuth
// @Tags shop
// @Description Input params for search shops
// @ModuleID shops
// @Accept  json
// @Produce  json
// @Param input query ShopInput true "params for search shops"
// @Success 200 {object} []domain.Shop
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/shops [get].
func (h *HandlerV1) FindShop(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Shop{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	shops, err := h.services.Shop.FindShop(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, shops)
}

func (h *HandlerV1) GetShopByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateShop(c *gin.Context) {

}

func (h *HandlerV1) DeleteShop(c *gin.Context) {

}
