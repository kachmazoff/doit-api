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

// @Summary Create suggestion
// @Security Auth
// @Tags suggestions
// @Description Создание нового предложения для участника
// @Accept json
// @Produce json
// @Param participantId path string true "Id участника"
// @Param input body model.Suggestion true "Модель предложения"
// @Success 200 {object} dto.IdResponse
// @Failure 400,403 {object} dto.MessageResponse
// @Router /participants/{participantId}/suggestions [post]
func (h *Controller) createSuggestion(c *gin.Context) {
	currentUser, _ := getUserId(c)
	participantId := c.Param("participantId")

	if !h.services.Participants.IsPublic(participantId) {
		c.AbortWithStatusJSON(http.StatusForbidden, createMessage("Вы не можете предлагать что-либо данному пользователю"))
		return
	}

	var input model.Suggestion
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage(err.Error()))
		return
	}

	input.AuthorId = currentUser
	input.ParticipantId = participantId

	id, err := h.services.Suggestions.Create(input)

	handleCreation(c, id, err)
}

// @Summary Get suggestions
// @Tags suggestions
// @Description Получение списка предложений для участника
// @Accept json
// @Produce json
// @Param participantId path string true "Id участника"
// @Success 200 {array} model.Suggestion
// @Failure 404 {object} dto.MessageResponse
// @Router /participants/{participantId}/suggestions [get]
func (h *Controller) getSuggestionsForParticipant(c *gin.Context) {
	participantId := c.Param("participantId")

	if participantId == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, createMessage("Participant is not defined"))
		return
	}
	suggestions, err := h.services.Suggestions.GetForParticipant(participantId)
	commonJSONResponse(c, suggestions, err)
}
