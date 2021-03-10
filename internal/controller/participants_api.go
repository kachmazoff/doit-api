package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) initParticipantsRoutes(api *gin.RouterGroup) {
	participants := api.Group("/participants")
	{
		participants.POST("/", h.userIdentity, h.createParticipant)
	}
}

func (h *Controller) createParticipant(c *gin.Context) {
	var input model.Participant
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	id, err := h.services.Participants.Create(input)
	handleCreation(c, id, err)
}
