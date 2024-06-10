package bearer

import (
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

type BearerAuthenticator struct {
	Auth map[string]interface{}
}

func (a *BearerAuthenticator) Authenticate(c *gin.Context) {

	username := a.Auth["username"].(string)
	password := a.Auth["passowrd"].(string)

	if username == "" || password == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": ErrCredentials})
		return
	}

	auth := username + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	c.Set("auth", encodedAuth)

	c.Next()
}
