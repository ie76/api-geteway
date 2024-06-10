package geolocation

import (
	"assignment/external"
	mainErrors "assignment/internal/errors"
	"assignment/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GeolocationService struct {
	Config *external.Config
}

func NewGeolocationService() *GeolocationService {
	conf, err := external.GetServiceConfig("geolocation")
	if err != nil {
		panic(mainErrors.New(err.Code))
	}

	return &GeolocationService{
		Config: conf,
	}
}

// GeolocationGet godoc
//
// @Summary geolocation
// @Description get your location by ip
// @Tags geolocation
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Router /geolocation [get]
func (s *GeolocationService) Do(c *gin.Context) (*http.Response, error) {
	endpoint := c.Param("path")

	var response *http.Response
	if c.Request.Method == "GET" {

		getKey := s.Config.AuthCredentials["key"]

		url := s.Config.BaseUrl + endpoint
		newUrl, err := utils.AddQueryParams(url, map[string]string{
			"key": getKey.(string),
		})

		if err != nil {
			return nil, err
		}

		resp, errGet := http.Get(newUrl)
		if errGet != nil {
			return nil, errGet
		}

		response = resp
	}

	return response, nil
}

func (s *GeolocationService) Authenticate(c *gin.Context) {
	auth := &GeolocationAuthenticator{
		Auth: s.Config.AuthCredentials,
	}
	auth.Authenticate(c)
}

func (s *GeolocationService) GetCacheDuration() int {
	return s.Config.CacheDuration
}
