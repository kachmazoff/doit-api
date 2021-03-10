package controller

import "github.com/gin-gonic/gin"

func (h *Controller) initFollowersRoutes(api *gin.RouterGroup) {
	api.POST("/follow", h.userIdentity, h.follow)
	api.POST("/unfollow", h.userIdentity, h.unfollow)
}

func (h *Controller) follow(c *gin.Context) {
	currentUser, _ := getUserId(c)
	otherUser := c.Query("userId")
	err := h.services.Followers.Subscribe(currentUser, otherUser)
	commonJSONResponse(c, createMessage("Вы успешно подписались"), err)
}

func (h *Controller) unfollow(c *gin.Context) {
	currentUser, _ := getUserId(c)
	otherUser := c.Query("userId")
	err := h.services.Followers.Unsubscribe(currentUser, otherUser)
	commonJSONResponse(c, createMessage("Вы успешно отписались"), err)
}
