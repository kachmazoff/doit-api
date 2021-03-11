package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) initChallengesRoutes(api *gin.RouterGroup) {
	courses := api.Group("/challenges")
	{
		courses.GET("/", h.getAllChallenges)
		courses.POST("/", h.userIdentity, h.createChallenge)
		courses.GET("/:challengeId/participants", h.getParticipantsByChallenge)
	}
}

// @Summary Create challenge
// @Security Auth
// @Tags challenges
// @Description Создание нового челленджа
// @Accept json
// @Produce json
// @Param input body model.Challenge true "Модель челленджа"
// @Success 200 {object} dto.IdResponse
// @Router /challenges [post]
func (h *Controller) createChallenge(c *gin.Context) {
	var input model.Challenge
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	currentUser, _ := getUserId(c)
	input.AuthorId = currentUser

	id, err := h.services.Challenges.Create(input)
	handleCreation(c, id, err)
}

// @Summary Get all challenges
// @Tags challenges
// @Description Получение списка челленджей
// @Accept json
// @Produce json
// @Success 200 {array} model.Challenge
// @Router /challenges [get]
func (h *Controller) getAllChallenges(c *gin.Context) {
	challenges, err := h.services.Challenges.GetAll()
	commonJSONResponse(c, challenges, err)
}

// @Summary Get challenge's participants
// @Tags challenges, participants
// @Description Получение списка участников в челлендже
// @Accept json
// @Produce json
// @Param challengeId path string true "Id челленджа"
// @Param status query string false "Статус участников"
// @Success 200 {array} model.Participant
// @Router /challenges/{challengeId}/participants [get]
func (h *Controller) getParticipantsByChallenge(c *gin.Context) {
	challengeId := c.Param("challengeId")

	status := c.Query("status")
	onlyActive := status == "active"
	participants, err := h.services.Participants.GetParticipantsInChallenge(challengeId, true, onlyActive)
	commonJSONResponse(c, participants, err)
}
