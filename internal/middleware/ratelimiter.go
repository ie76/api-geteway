package middleware

import (
	"assignment/internal/errors"
	"assignment/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimiter(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in request context"})
		return
	}

	user, errGet := services.NewUserService().GetUserById(userID.(int))
	if errGet != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errors.ErrUserNotFound})
		return
	}

	errDeduction := services.NewUserService().DeductCredits(user.ID)
	if errDeduction != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInsufficientCredits})
		return
	}

	c.Next()

}
