package controller

import "github.com/gin-gonic/gin"

func (h *Controller) initFollowersRoutes(api *gin.RouterGroup) {
	api.POST("/follow", h.userIdentity, h.follow)
	api.POST("/unfollow", h.userIdentity, h.unfollow)
}

// @Summary Follow user
// @Security Auth
// @Tags followers
// @Description Подписка на пользователя
// @Accept json
// @Produce json
// @Param userId query string true "Id пользователя"
// @Success 200 {object} dto.MessageResponse
// @Failure 403 {object} dto.MessageResponse
// @Router /follow [post]
func (h *Controller) follow(c *gin.Context) {
	currentUser, _ := getUserId(c)
	otherUser := c.Query("userId")
	err := h.services.Followers.Subscribe(currentUser, otherUser)
	commonJSONResponse(c, createMessage("Вы успешно подписались"), err)
}

// @Summary Unfollow user
// @Security Auth
// @Tags followers
// @Description Отписка от пользователя
// @Accept json
// @Produce json
// @Param userId query string true "Id пользователя"
// @Success 200 {object} dto.MessageResponse
// @Failure 403 {object} dto.MessageResponse
// @Router /unfollow [post]
func (h *Controller) unfollow(c *gin.Context) {
	currentUser, _ := getUserId(c)
	otherUser := c.Query("userId")
	err := h.services.Followers.Unsubscribe(currentUser, otherUser)
	commonJSONResponse(c, createMessage("Вы успешно отписались"), err)
}
