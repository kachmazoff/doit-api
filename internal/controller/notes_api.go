package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) initNotesRoutes(api *gin.RouterGroup) {
	notes := api.Group("/notes")
	{
		notes.POST("/", h.userIdentity, h.createNote)
		notes.GET("/", h.optionalUserIdentity, h.getParticipantNotes)
	}
}

func (h *Controller) createNote(c *gin.Context) {
	currentUser, _ := getUserId(c)
	participantId := c.Param("participantId")

	if !h.services.Participants.HasRootAccess(participantId, currentUser) {
		c.AbortWithStatusJSON(http.StatusForbidden, "Вы не можете добавлять записи от лица данного участника")
		return
	}

	var input model.Note
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	input.AuthorId = currentUser
	input.ParticipantId = participantId

	id, err := h.services.Notes.Create(input)
	handleCreation(c, id, err)
}

func (h *Controller) getParticipantNotes(c *gin.Context) {
	participantId := c.Param("participantId")
	hasRootAccess := false

	currentUser, err := getUserId(c)
	if err == nil {
		hasRootAccess = h.services.Participants.HasRootAccess(participantId, currentUser)
	}

	var isPublic bool
	if !hasRootAccess {
		isPublic = h.services.Participants.IsPublic(participantId)
	}

	if !hasRootAccess && !isPublic {
		c.AbortWithStatusJSON(http.StatusForbidden, "Вы не можете просматривать записи данного участника")
		return
	}

	needAnonymize := isPublic && !hasRootAccess
	notes, err := h.services.Notes.GetNotesOfParticipant(participantId, needAnonymize)

	commonJSONResponse(c, notes, err)
}
