package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/docs"
	"github.com/mikalai2006/go-template-api/internal/config"
	v1 "github.com/mikalai2006/go-template-api/internal/handler/v1"
	"github.com/mikalai2006/go-template-api/internal/middleware"
	"github.com/mikalai2006/go-template-api/internal/service"
	"github.com/mikalai2006/go-template-api/pkg/app"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Services
	oauth    config.OauthConfig
	i18n     config.I18nConfig
}

func NewHandler(services *service.Services, oauth *config.OauthConfig, i18n *config.I18nConfig) *Handler {
	return &Handler{
		services: services,
		oauth:    *oauth,
		i18n:     *i18n,
	}
}

func (h *Handler) InitRoutes(cfg *config.Config) *gin.Engine {
	// appG := app.Gin{C: *gin.Context}
	router := gin.New()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.Cors,
		// middleware.JSONAppErrorReporter(),
	)
	// add swagger route
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if cfg.Environment != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.HTTP.Host
	}
	if cfg.Environment != config.Prod {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router.GET("/", func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusOK, "API")
	})

	// create session
	// store := cookie.NewStore([]byte(os.Getenv("secret")))
	// router.Use(sessions.Sessions("mysession", store))

	h.initAPI(router)

	router.NoRoute(func(c *gin.Context) {
		appG := app.Gin{C: c}
		// c.AbortWithError(http.StatusNotFound, errors.New("page not found"))
		appG.ResponseError(http.StatusNotFound, errors.New("page not found"), nil)
		// .SetMeta(gin.H{
		// 	"code": http.StatusNotFound,
		// 	"status": "error",
		// 	"message": "hello",
		// })
	})
	router.Static("/images", "./public")
	router.Static("/css", "./public/css")
	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.services, &h.oauth, &h.i18n)
	api := router.Group("/api")
	handlerV1.Init(api)
}
