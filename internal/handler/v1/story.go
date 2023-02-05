package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/domain"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/utils"
	"github.com/mikalai2006/go-template-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
)

func (h HandlerV1) RegisterStory(router *gin.RouterGroup) {
	story := router.Group("/story")
	story.POST("/:id", middleware.SetUserIdentity, h.publishStory)
	story.GET("/*slug", h.getStory)
}

func (h HandlerV1) publishStory(c *gin.Context) {
	appG := app.Gin{C: c}

	var result domain.Story

	var input domain.StoryInputData
	err := c.Bind(&input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	input.UserID = userID

	id := c.Param("id")

	input.PageID = id

	result, err = h.services.Story.PublishStory(id, input)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h HandlerV1) getStory(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.StoryInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// add slug in params.
	slug := c.Param("slug")
	slug = strings.TrimPrefix(slug, "/")

	if slug != "" {
		params.Filter.(bson.M)["slug_full"] = slug
	}

	// fmt.Println("params ", params)

	document, err := h.services.Story.GetStory(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// if document.PageID.Hex() != "" {
	// 	filePath := fmt.Sprintf("./public/css/p_%v.css", document.PageID.Hex())
	// 	f, err := ioutil.ReadFile(filePath)
	// 	if err != nil {
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	}
	// 	document.CSS = string(f)
	// }

	c.JSON(http.StatusOK, document)
}
