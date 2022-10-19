package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/go-template-api/pkg/app"
	"github.com/mikalai2006/go-template-api/pkg/auths"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userRoles           = "roles"
)

func SetUserIdentity(c *gin.Context) {
	appG := app.Gin{C: c}

	header := c.GetHeader(authorizationHeader)

	if header == "" {
		// c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("empty auth header"))
		appG.Response(http.StatusUnauthorized, errors.New("empty auth header"), nil)
		return
	}

	headerParts := strings.Split(header, " ")
	countParts := 2
	if len(headerParts) != countParts {
		// c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
		appG.Response(http.StatusUnauthorized, errors.New("invalid auth header"), nil)
		return
	}

	if headerParts[1] == "" {
		// c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
		appG.Response(http.StatusUnauthorized, errors.New("invalid auth header"), nil)
		return
	}

	// parse token
	// userId, err := h.services.Authorization.ParseToken(headerParts[1])
	// if err != nil {
	// 	newErrorResponse(c, http.StatusUnauthorized, err.Error())
	// 	return
	// }
	tokenManager, err := auths.NewManager(os.Getenv("SIGNING_KEY"))
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.Response(http.StatusUnauthorized, err, nil)
		return
	}

	claims, err := tokenManager.Parse(headerParts[1])
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.Response(http.StatusUnauthorized, err, nil)
		return
	}
	c.Set(userCtx, claims.Subject)
	c.Set(userRoles, claims.Roles)
	// session := sessions.Default(c)
	// user := session.Get(userkey)
	// if user == nil {
	// 	// Abort the request with the appropriate error code
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	// logrus.Printf("user session= %s", user)
	// // Continue down the chain to handler etc
	// c.Next()
}

func GetUserID(c *gin.Context) (string, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return "", errors.New("user not found")
	}

	idInt, ok := id.(string)
	if !ok {
		return "", errors.New("user not found2")
	}

	return idInt, nil
}

func GetRoles(c *gin.Context) ([]string, error) {
	roles, ok := c.Get(userRoles)
	if !ok {
		return nil, errors.New("roles not found")
	}
	return roles.([]string), nil
}
