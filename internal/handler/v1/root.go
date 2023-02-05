package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/service"
)

type HandlerV1 struct {
	services    *service.Services
	oauth       config.OauthConfig
	i18n        config.I18nConfig
	imageConfig config.IImageConfig
}

func NewHandler(services *service.Services, oauth *config.OauthConfig, i18n *config.I18nConfig, imageConfig *config.IImageConfig) *HandlerV1 {
	return &HandlerV1{
		services:    services,
		oauth:       *oauth,
		i18n:        *i18n,
		imageConfig: *imageConfig,
	}
}

func (h *HandlerV1) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{

		h.registerAuth(v1)
		oauth := v1.Group("/oauth")
		h.registerVkOAuth(oauth)
		h.registerGoogleOAuth(oauth)

		h.registerShop(v1)
		h.RegisterUser(v1)
		h.RegisterPage(v1)
		h.RegisterComponent(v1)
		h.RegisterComponentGroup(v1)
		h.RegisterComponentPreset(v1)
		h.RegisterApp(v1)
		h.RegisterProduct(v1)
		h.RegisterImage(v1)
		h.RegisterSpace(v1)
		h.RegisterPlugin(v1)
		h.RegisterPartner(v1)
		h.RegisterStory(v1)

		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "v1",
			})
		})
	}
}
