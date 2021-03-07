package controller

import (
	"net/http"
	"time"

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
		api.POST("/auth/login", h.getToken)
		api.POST("/auth/registration", h.registerUser)
		api.POST("/auth/activate", h.activateAccount)
		api.GET("/user", h.getUser)
		h.initChallengesRoutes(api)
		h.initTimelineRoutes(api)
	}

	return router
}

// TODO: перенести в другой файл
func (h *Controller) getToken(c *gin.Context) {
	var input map[string]interface{}
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	token, err := h.tokenManager.NewJWT(input["id"].(string), time.Hour)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
