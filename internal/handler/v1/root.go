package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/internal/config"
	"github.com/mikalai2006/go-template-api/internal/service"
)

type HandlerV1 struct {
	services *service.Services
	oauth    config.OauthConfig
	i18n     config.I18nConfig
}

func NewHandler(services *service.Services, oauth *config.OauthConfig, i18n *config.I18nConfig) *HandlerV1 {
	return &HandlerV1{
		services: services,
		oauth:    *oauth,
		i18n:     *i18n,
	}
}

func (h *HandlerV1) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		auth := v1.Group("/auth")
		h.registerAuth(auth)

		h.registerShop(v1)
		h.RegisterUser(v1)
		h.RegisterPage(v1)
		h.RegisterComponent(v1)
		h.RegisterApp(v1)
		h.RegisterProduct(v1)

		oauth := v1.Group("/oauth")
		h.registerVkOAuth(oauth)
		h.registerGoogleOAuth(oauth)

		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "v1",
			})
		})
	}
}
