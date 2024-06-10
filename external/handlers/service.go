package handlers

import (
	"assignment/external"
	"assignment/internal/cache"
	"assignment/internal/errors"
	b64 "encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//TODO: SEND TO EXTERNAL

func HandleRequest(service external.ExternalService, cache *cache.RedisCache, serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		key := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s_%s", serviceName, c.Request.URL.Path)))
		_, ok := c.GetQuery("force_refresh")

		if !ok {
			if data, err := cache.Get(c, key); err == nil {
				c.Data(http.StatusOK, "application/json", data)
				return
			}
		}

		resp, body, err := callExternalService(service, c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": err.Code, "message": err.Message, "at": err.Time})
			return
		}

		errCaching := cache.Set(c.Request.Context(), key, body, time.Duration(service.GetCacheDuration())*time.Second)
		if errCaching != nil {
			tError := errors.New(errors.ErrCacheError)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": tError.Code, "message": tError.Message, "at": tError.Time})
			return
		}

		c.Data(resp.StatusCode, "application/json", body)
	}
}

func callExternalService(service external.ExternalService, c *gin.Context) (*http.Response, []byte, *errors.Error) {
	// Call the external service
	resp, err := service.Do(c)
	if err != nil {
		return nil, nil, errors.New(errors.ErrServiceHTTPRequest)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, nil, errors.New(errors.ErrServiceReadBody)
	}

	return resp, body, nil
}
