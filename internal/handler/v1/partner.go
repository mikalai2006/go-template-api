package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *HandlerV1) RegisterPartner(router *gin.RouterGroup) {
	Partner := router.Group("/partner")
	Partner.GET("", h.findPartner)
	Partner.POST("", middleware.SetUserIdentity, h.createPartner)
	Partner.GET("/:id", h.getPartner)
	Partner.PATCH("/:id", middleware.SetUserIdentity, h.updatePartner)
	Partner.DELETE("/:id", middleware.SetUserIdentity, h.deletePartner)
}

func (h *HandlerV1) createPartner(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input = &domain.PartnerInput{}
	if er := c.Bind(input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	document, err := h.services.Partner.CreatePartner(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	// fmt.Println("input", input)
	// input.UserID = userID
	// // var image domain.Image

	var imageInput = &domain.ImageInput{}
	imageInput.Service = "partner"
	imageInput.ServiceID = document.ID.Hex()
	imageInput.UserID = userID
	imageInput.Dir = "partner"

	paths, err := utils.UploadResizeMultipleFile(c, imageInput, "upload", &h.imageConfig)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
	}

	var result []domain.Image
	for i := range paths {
		imageInput.Path = paths[i].Path
		imageInput.Ext = paths[i].Ext
		image, err := h.services.Image.CreateImage(userID, imageInput)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, image)
	}
	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) getPartner(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	document, err := h.services.Partner.GetPartner(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) findPartner(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.PartnerInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	documents, err := h.services.Partner.FindPartner(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	for i := range documents.Data {
		params := domain.RequestParams{
			Filter: bson.D{{"service", "partner"}, {"service_id", documents.Data[i].ID.Hex()}},
		}
		images, err := h.services.Image.FindImage(params)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		documents.Data[i].Images = images.Data

		user, err := h.services.User.GetUser(documents.Data[i].UserID.Hex())
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		documents.Data[i].User = user
	}

	c.JSON(http.StatusOK, documents)
}

func (h *HandlerV1) updatePartner(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input = &domain.PartnerInput{}
	if er := c.Bind(input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println("data", input)

	document, err := h.services.Partner.UpdatePartner(id, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	var imageInput = &domain.ImageInput{}
	imageInput.Service = "partner"
	imageInput.ServiceID = document.ID.Hex()
	imageInput.UserID = userID
	imageInput.Dir = "partner"

	paths, err := utils.UploadResizeMultipleFile(c, imageInput, "upload", &h.imageConfig)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
	}

	var result []domain.Image
	for i := range paths {
		imageInput.Path = paths[i].Path
		imageInput.Ext = paths[i].Ext
		image, err := h.services.Image.CreateImage(userID, imageInput)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, image)
	}

	// var input domain.PartnerInput
	// data, err := utils.BindAndValid(c, &input)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// // fmt.Println(data)

	// document, err := h.services.Partner.UpdatePartner(id, &data)
	// if err != nil {
	// 	appG.ResponseError(http.StatusInternalServerError, err, nil)
	// 	return
	// }

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) deletePartner(c *gin.Context) {
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

	document, err := h.services.Partner.DeletePartner(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}
