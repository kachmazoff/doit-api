package controller

import "github.com/gin-gonic/gin"

func (h *Controller) initParticipantsRoutes(api *gin.RouterGroup) {
	participant := api.Group("/participants/:participantId")
	{
		h.initNotesRoutes(participant)
		h.initSuggestionsRoutes(participant)
	}
}
