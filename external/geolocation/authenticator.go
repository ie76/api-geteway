package geolocation

import (
	"github.com/gin-gonic/gin"
)

type GeolocationAuthenticator struct {
	Auth map[string]interface{}
}

func (a *GeolocationAuthenticator) Authenticate(c *gin.Context) {
	apiKey, ok := a.Auth["key"]

	if !ok || apiKey == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": ErrCredentials})
		return
	}
	c.Next()
}
