package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("/", h.getUser)
		users.GET("/:username/participants", h.optionalUserIdentity, h.getUserParticipations)
		users.GET("/:username/followees", h.getFollowees)
		users.GET("/:username/followers", h.getFollowers)
	}
}

func (h *Controller) getUser(c *gin.Context) {
	username := c.Query("username")

	if username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Bad request"))
		return
	}

	user, err := h.services.Users.GetByUsername(username)
	if handleUserFindError(c, err) {
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
	if handleUserFindError(c, err) {
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
	commonJSONResponse(c, participations, err)
}

func (h *Controller) getFollowees(c *gin.Context) {
	h.helperGetFollowX(c, h.services.Followers.GetFollowees)
}

func (h *Controller) getFollowers(c *gin.Context) {
	h.helperGetFollowX(c, h.services.Followers.GetFollowers)
}

func (h *Controller) helperGetFollowX(c *gin.Context, get func(userId string) ([]model.User, error)) {
	username := c.Param("username")
	userId, err := h.services.Users.GetIdByUsername(username)

	// TODO: проверить, возмонжно ли userId == ""
	if err != nil || userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Bad request. Мы не нашли такого пользователя"))
		return
	}

	followers, err := get(userId)
	commonJSONResponse(c, followers, err)
}
