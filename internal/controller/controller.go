package controller

import (
	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/kachmazoff/doit-api/docs"
	"github.com/kachmazoff/doit-api/internal/auth"
	"github.com/kachmazoff/doit-api/internal/service"
)

type Controller struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewController(services *service.Services, tokenManager auth.TokenManager) *Controller {
	return &Controller{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Controller) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.GET("/user", h.getUser)
		h.initAccountRoutes(api)
		h.initChallengesRoutes(api)
		h.initTimelineRoutes(api)
	}

	return router
}
