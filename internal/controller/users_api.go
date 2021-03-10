package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("/", h.getUser)
		users.GET("/:username/participants", h.optionalUserIdentity, h.getUserParticipations)
	}
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

func (h *Controller) getUserParticipations(c *gin.Context) {
	username := c.Param("username")
	status := c.Query("status")

	// TODO: check. Может ли такое быть вообще?
	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Bad request. Username not defined"))
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

	currentUser, err := getUserId(c)
	var onlyPublic bool

	if err != nil {
		onlyPublic = true
	} else {
		onlyPublic = currentUser != user.Id
	}

	participations, err := h.services.Participants.GetParticipationsOfUser(user.Id, onlyPublic, status == "active")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage(err.Error()))
		return
	}
	c.JSON(http.StatusOK, participations)
}
