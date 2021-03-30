package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Controller) initDebugRoutes(api *gin.RouterGroup) {
	debug := api.Group("/debug")
	{
		debug.GET("/genPass", h.genPass)
	}
}

func (h *Controller) genPass(c *gin.Context) {
	password := c.Query("password")

	if password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Bad request"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	commonJSONResponse(c, string(hashedPassword), err)
}
