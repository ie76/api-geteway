package middleware

import (
	"assignment/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	// Get the cookie off the request
	token := c.Request.Header.Get("Authorization")
	// gToken := token[len("Bearer "):]
	claim, err := utils.DecodeToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    err.Code,
			"at":      err.Time,
			"message": err.Message,
		})

		return
	}

	ok := claim["user_id"]
	c.Set("user_id", int(ok.(float64)))

	c.Next()
}
