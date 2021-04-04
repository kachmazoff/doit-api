package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) initParticipantsRoutes(api *gin.RouterGroup) {
	participant := api.Group("/participants/:participantId")
	{
		participant.GET("/", h.optionalUserIdentity, h.getParticipant)
		h.initNotesRoutes(participant)
		h.initSuggestionsRoutes(participant)
	}
}

// @Summary Get participant info
// @Tags participants
// @Description Get participant info by id
// @Accept  json
// @Produce  json
// @Param participantId path string true "id дневника"
// @Success 200 {object} model.Participant
// @Router /participants/{participantId} [get]
func (h *Controller) getParticipant(c *gin.Context) {
	participantId := c.Param("participantId")
	participant, err := h.services.Participants.GetByIdUNSAFE(participantId)

	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, createMessage("This participant does not exist"))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage("We are already working on it"))
		return
	}

	currentUserId, err := getUserId(c)

	if participant.VisibleType == "private" && participant.UserId != currentUserId {
		c.AbortWithStatusJSON(http.StatusNotFound, createMessage("This participant does not exist"))
		return
	}

	h.services.Participants.Anonymize(&participant, currentUserId)

	commonJSONResponse(c, participant, nil)
}
