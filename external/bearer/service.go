package bearer

import (
	"assignment/external"
	mainErrors "assignment/internal/errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BearerService struct {
	Config *external.Config
	Header string
}

func NewBearerService() *BearerService {
	conf, err := external.GetServiceConfig("bearer")
	if err != nil {
		panic(mainErrors.New(err.Code))
	}

	return &BearerService{
		Config: conf,
	}
}

// BearerSrvice  godoc
//
// @Summary Basic Auth Service Test
// @Description test a basic auth connexion
// @Tags Basic Auth
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Router /basic-auth [get]
func (s *BearerService) Do(c *gin.Context) (*http.Response, error) {
	var response *http.Response
	if c.Request.Method == "GET" {

		req, err := http.NewRequest("GET", s.Config.BaseUrl, nil)
		if err != nil {
			return response, nil
		}

		bearerHeader, ok := c.Get("auth")
		if !ok {
			return response, errors.New("failedTOPArse")
		}

		req.Header.Add("Authorization", "Bearer "+bearerHeader.(string))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		response = resp
	}

	return response, nil
}

func (s *BearerService) Authenticate(c *gin.Context) {
	auth := &BearerAuthenticator{
		Auth: s.Config.AuthCredentials,
	}
	auth.Authenticate(c)
}

func (s *BearerService) GetCacheDuration() int {
	return s.Config.CacheDuration
}
