package external

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExternalService interface {
	Do(c *gin.Context) (*http.Response, error)
	Authenticate(c *gin.Context)
	GetCacheDuration() int
}

var ServiceRegistry = map[string]ExternalService{}

func RegisterService(name string, service ExternalService) {
	ServiceRegistry[name] = service
}
