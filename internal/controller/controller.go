package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/kachmazoff/doit-api/docs"
	"github.com/kachmazoff/doit-api/internal/service"
)

type Controller struct {
	services *service.Services
}

func NewController(services *service.Services) *Controller {
	return &Controller{services: services}
}

func (h *Controller) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.GET("/hello", h.helloWorld)
		api.POST("/auth/registration", h.registerUser)
		api.POST("/auth/activate", h.activateAccount)
		api.GET("/user", h.getUser)
	}

	return router
}

func (h *Controller) helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, "hello")
}
