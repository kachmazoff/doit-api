package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) registerUser(c *gin.Context) {
	var input model.User
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	id, err := h.services.Users.Create(input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Controller) getUser(c *gin.Context) {
	username := c.Query("username")

	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Bad request"))
		return
	}

	user, err := h.services.Users.GetByUsername(username)

	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, createMessage("This user does not exist"))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Controller) activateAccount(c *gin.Context) {
	userId := c.Query("id")

	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Bad request"))
		return
	}

	err := h.services.Users.ConfirmAccount(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, createMessage("Account activated successfully"))
}

func createMessage(message string) map[string]interface{} {
	return map[string]interface{}{
		"message": "Account activated successfully",
	}
}
