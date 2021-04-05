package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) initChallengesRoutes(api *gin.RouterGroup) {
	challenges := api.Group("/challenges")
	{
		challenges.GET("/", h.optionalUserIdentity, h.getAllChallenges)
		// challenges.GET("/public", h.getAllPublicChallenges)
		// challenges.GET("/own", h.userIdentity, h.getAllOwnChallenges)
		challenges.POST("/", h.userIdentity, h.createChallenge)

		challenge := challenges.Group("/:challengeId")
		{
			challenge.GET("/", h.optionalUserIdentity, h.getChallengeById)
			challenge.GET("/participants", h.getParticipantsByChallenge)
			challenge.POST("/participants", h.userIdentity, h.createParticipant)
		}
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
	currentUser, _ := getUserId(c)
	challengesType := c.Query("type")
	var challenges []model.Challenge
	var err error

	if challengesType == "" {
		challenges, err = h.services.Challenges.GetAll()
	} else if challengesType == "public" {
		challenges, err = h.services.Challenges.GetAllPublic()
	} else if challengesType == "own" {
		if currentUser == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, createMessage("Необходимо авторизоваться"))
			return
		}
		challenges, err = h.services.Challenges.GetAllOwn(currentUser)
	}
	if challenges != nil && len(challenges) > 0 && currentUser != "" {
		h.services.Challenges.EnrichAllWithUserParticipant(&challenges, currentUser)
	}
	commonJSONResponse(c, challenges, err)
}

// @Summary Get challenge's info
// @Tags challenges
// @Description Get challenge's info by id
// @Accept  json
// @Produce  json
// @Param challengeId path string true "id челленджа"
// @Success 200 {object} model.Challenge
// @Router /challenges/{challengeId} [get]
func (h *Controller) getChallengeById(c *gin.Context) {
	currentUser, _ := getUserId(c)
	challengeId := c.Param("challengeId")
	challenge, err := h.services.Challenges.GetById(challengeId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Мы не нашли такого челленджа"))
		return
	}

	if currentUser != "" {
		h.services.Challenges.EnrichWithUserParticipant(&challenge, currentUser)
	}

	commonJSONResponse(c, challenge, err)
}

// @Summary Get all public challenges
// @Tags challenges
// @Description Получение списка публичных челленджей
// @Accept json
// @Produce json
// @Success 200 {array} model.Challenge
// @Router /challenges/public [get]
func (h *Controller) getAllPublicChallenges(c *gin.Context) {
	challenges, err := h.services.Challenges.GetAllPublic()
	commonJSONResponse(c, challenges, err)
}

// @Summary Get all own challenges
// @Security Auth
// @Tags challenges
// @Description Получение списка личных челленджей
// @Accept json
// @Produce json
// @Success 200 {array} model.Challenge
// @Router /challenges/own [get]
func (h *Controller) getAllOwnChallenges(c *gin.Context) {
	currentUser, _ := getUserId(c)
	challenges, err := h.services.Challenges.GetAllOwn(currentUser)
	commonJSONResponse(c, challenges, err)
}

// @Summary Get challenge's participants
// @Tags participants
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

// @Summary Create participant
// @Security Auth
// @Tags participants
// @Description Создание нового участника (регистрация в челлендже в качестве участника)
// @Accept json
// @Produce json
// @Param challengeId path string true "Id челленджа"
// @Param input body model.Participant true "Модель участника"
// @Success 200 {object} dto.IdResponse
// @Failure 400,403 {object} dto.MessageResponse
// @Router /challenges/{challengeId}/participants [post]
func (h *Controller) createParticipant(c *gin.Context) {
	challengeId := c.Param("challengeId")

	var input model.Participant
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	currentUser, _ := getUserId(c)
	input.UserId = currentUser
	input.ChallengeId = challengeId

	// TODO: test
	id, err := h.services.Participants.Create(input)
	handleCreation(c, id, err)
}
