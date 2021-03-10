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
		courses.GET("/:challengeId", h.getAllChallenges)
	}
}

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

func (h *Controller) getAllChallenges(c *gin.Context) {
	challenges, err := h.services.Challenges.GetAll()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, challenges)
}

func (h *Controller) getParticipantsByChallenge(c *gin.Context) {
	challengeId := c.Param("challengeId")

	status := c.Query("status")

	participants, err := h.services.Participants.GetParticipantsInChallenge(challengeId, true, status == "active")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage(err.Error()))
		return
	}
	c.JSON(http.StatusOK, participants)
}
