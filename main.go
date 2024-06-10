package main

import (
	"assignment/config"
	"assignment/database"
	"assignment/external"
	"assignment/internal/auth"
	"assignment/internal/cache"
	"assignment/internal/errors"
	"assignment/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "assignment/docs"
	_ "assignment/external/bearer"
	_ "assignment/external/geolocation"
	extHandlers "assignment/external/handlers"
)

// @title Assignement API
// @version 1.0
// @description Assignement for Issam Elyazidi.
// @host localhost:8080
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
func main() {

	err := config.NewConfig()
	if err != nil {
		panic(errors.New(err.Code))
	}

	if err := database.Init(); err != nil {
		panic(errors.New(errors.ErrDatabaseError))
	}
	defer database.GetDB().Close()

	database.Migrate()

	cacheProvider, err := cache.NewCacheProvider()
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/login", auth.Login)
	r.POST("/register", auth.Register)

	for name, service := range external.ServiceRegistry {
		serviceRoute := r.Group("/"+name).Use(middleware.Authenticate, middleware.RateLimiter)
		{
			serviceRoute.Use(service.Authenticate)
			serviceRoute.Any("/*path", extHandlers.HandleRequest(service, cacheProvider, name))
		}
	}

	r.Run(":8080")
}
