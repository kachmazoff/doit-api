package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) createChallenge(c *gin.Context) {
	var input model.Challenge
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	id, err := h.services.Challenges.Create(input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Controller) getChallenges(c *gin.Context) {
	challenges, err := h.services.Challenges.GetAll()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, challenges)
}
