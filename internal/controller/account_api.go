package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) initAccountRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/registration", h.registerUser)
		auth.POST("/login", h.getToken)
		auth.POST("/activate", h.activateAccount)
	}
}

func (h *Controller) registerUser(c *gin.Context) {
	var input model.User
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	id, err := h.services.Users.Create(input)
	handleCreation(c, id, err)
}

func (h *Controller) activateAccount(c *gin.Context) {
	userId := c.Query("id")

	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Bad request"))
		return
	}

	err := h.services.Users.ConfirmAccount(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, createMessage("Account activated successfully"))
}

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
