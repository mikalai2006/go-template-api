package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/app"
)

func (h *HandlerV1) RegisterProduct(router *gin.RouterGroup) {
	product := router.Group("/product")
	product.GET("", h.findProduct)
	product.POST("", middleware.SetUserIdentity, h.createProduct)
	product.GET("/:id", h.getProduct)
	product.PATCH("/:id", middleware.SetUserIdentity, h.updateProduct)
	product.DELETE("/:id", middleware.SetUserIdentity, h.deleteProduct)
}

func (h *HandlerV1) createProduct(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var a map[string]interface{}
	if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	data, er := utils.BindJSON[domain.ProductInput](a)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// var input *domain.ProductInput
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.Response(http.StatusBadRequest, er, nil)
	// 	return
	// }

	document, err := h.services.Product.CreateProduct(userID, &data)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) getProduct(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	document, err := h.services.Product.GetProduct(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) findProduct(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.ProductInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	documents, err := h.services.Product.FindProduct(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, documents)
}

func (h *HandlerV1) updateProduct(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	var input domain.ProductInput
	data, err := utils.BindAndValid(c, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	// fmt.Println(data)

	document, err := h.services.Product.UpdateProduct(id, &data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) deleteProduct(c *gin.Context) {
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

	document, err := h.services.Product.DeleteProduct(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}
