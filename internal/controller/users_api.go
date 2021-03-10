package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) initUsersRoutes(api *gin.RouterGroup) {
	users := api.Group("/users/:username")
	{
		users.GET("/", h.getUser)
		users.GET("/participants", h.optionalUserIdentity, h.getUserParticipations)
		users.GET("/followees", h.getFollowees)
		users.GET("/followers", h.getFollowers)
	}
}

// @Summary Get user info
// @Tags users
// @Description Get user info by username
// @Accept  json
// @Produce  json
// @Param username path string true "username пользователя"
// @Success 200 {object} model.User
// @Router /users/{username} [get]
func (h *Controller) getUser(c *gin.Context) {
	username := c.Param("username")

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

// @Summary Get user's participations
// @Security Auth
// @Tags users, participants
// @Description Get participations of user by username
// @Accept  json
// @Produce  json
// @Param username path string true "username пользователя"
// @Param status query string false "status для фильтраций"
// @Success 200 {array} model.Participant
// @Router /users/{username}/participants [get]
func (h *Controller) getUserParticipations(c *gin.Context) {
	username := c.Param("username")
	status := c.Query("status")

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

// @Summary Get user's followees
// @Tags users, followers
// @Description Получение списка пользователей, на которых он подписан
// @Accept json
// @Produce json
// @Param username path string true "username пользователя"
// @Success 200 {array} model.User
// @Router /users/{username}/followees [get]
func (h *Controller) getFollowees(c *gin.Context) {
	h.helperGetFollowX(c, h.services.Followers.GetFollowees)
}

// @Summary Get user's followers
// @Tags users, followers
// @Description Получение списка пользователей, которые подписанны на данного пользователя
// @Accept json
// @Produce json
// @Param username path string true "username пользователя"
// @Success 200 {array} model.User
// @Router /users/{username}/followers [get]
func (h *Controller) getFollowers(c *gin.Context) {
	h.helperGetFollowX(c, h.services.Followers.GetFollowers)
}

func (h *Controller) helperGetFollowX(c *gin.Context, get func(userId string) ([]model.User, error)) {
	username := c.Param("username")
	userId, err := h.services.Users.GetIdByUsername(username)

	if err != nil || userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Bad request. Мы не нашли такого пользователя"))
		return
	}

	followers, err := get(userId)
	commonJSONResponse(c, followers, err)
}
