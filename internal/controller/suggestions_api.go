package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) initSuggestionsRoutes(api *gin.RouterGroup) {
	suggestions := api.Group("/suggestions")
	{
		suggestions.GET("/", h.getSuggestionsForParticipant)
		suggestions.POST("/", h.userIdentity, h.createSuggestion)
	}
}

func (h *Controller) createSuggestion(c *gin.Context) {
	currentUser, _ := getUserId(c)
	participantId := c.Param("participantId")

	if !h.services.Participants.IsPublic(participantId) {
		c.AbortWithStatusJSON(http.StatusForbidden, "Вы не можете предлагать что-либо данному пользователю")
		return
	}

	var input model.Suggestion
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	input.AuthorId = currentUser
	input.ParticipantId = participantId

	id, err := h.services.Suggestions.Create(input)

	handleCreation(c, id, err)
}

func (h *Controller) getSuggestionsForParticipant(c *gin.Context) {
	participantId := c.Param("participantId")
	// TODO: удалть проверку?
	if participantId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Participant is not defined"))
		return
	}
	suggestions, err := h.services.Suggestions.GetForParticipant(participantId)
	commonJSONResponse(c, suggestions, err)
}
