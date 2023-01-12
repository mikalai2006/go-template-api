package v1

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/app"
)

// func init() {
// 	if _, err := os.Stat("public/single"); errors.Is(err, os.ErrNotExist) {
// 		err := os.MkdirAll("public/single", os.ModePerm)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// 	if _, err := os.Stat("public/multiple"); errors.Is(err, os.ErrNotExist) {
// 		err := os.MkdirAll("public/multiple", os.ModePerm)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// }

func (h *HandlerV1) RegisterImage(router *gin.RouterGroup) {
	route := router.Group("/image")
	route.GET("/:id", h.getImage)
	route.GET("/find", h.findImage)
	route.POST("/", middleware.SetUserIdentity, h.createImage)
	route.DELETE("/:id", middleware.SetUserIdentity, h.deleteImage)
}

func (h *HandlerV1) getImage(c *gin.Context) {
	appG := app.Gin{C: c}
	id := c.Param("id")

	image, err := h.services.Image.GetImage(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, image)
}

func (h *HandlerV1) findImage(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Image{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	images, err := h.services.Image.FindImage(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, images)
}

func (h *HandlerV1) createImage(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	appG := app.Gin{C: c}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input = &domain.ImageInput{}
	if er := c.Bind(input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println("input", input)
	input.UserID = userID
	// var image domain.Image

	paths, err := utils.UploadResizeMultipleFile(c, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
	}
	input.Path = paths[0]
	image, err := h.services.Image.CreateImage(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	c.JSON(http.StatusOK, image)
}

func (h *HandlerV1) deleteImage(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	imageForRemove, err := h.services.Image.GetImage(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if imageForRemove.Service == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("not found item for remove"), nil)
		return
	} else {
		pathRemove := fmt.Sprintf("public/%s/%s/%s", imageForRemove.UserID.Hex(), imageForRemove.Service, imageForRemove.Path)
		err := os.Remove(pathRemove)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
		}
		pathRemove = fmt.Sprintf("public/%s/%s/xs-%s", imageForRemove.UserID.Hex(), imageForRemove.Service, imageForRemove.Path)
		err = os.Remove(pathRemove)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
		}
		pathRemove = fmt.Sprintf("public/%s/%s/lg-%s", imageForRemove.UserID.Hex(), imageForRemove.Service, imageForRemove.Path)
		err = os.Remove(pathRemove)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
		}
	}

	image, err := h.services.Image.DeleteImage(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, image)
}
