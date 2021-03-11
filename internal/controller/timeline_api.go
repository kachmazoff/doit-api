package controller

import "github.com/gin-gonic/gin"

func (h *Controller) initTimelineRoutes(api *gin.RouterGroup) {
	timeline := api.Group("/timeline")
	{
		timeline.GET("/", h.getCommonTimeline)
		timeline.GET("/own", h.userIdentity, h.getOwnTimeline)
		timeline.GET("/personalized", h.userIdentity, h.getPersonalizedTimeline)
	}
}

// @Summary Get common timeline
// @Tags timeline
// @Description Получение общего таймлайна
// @Accept json
// @Produce json
// @Success 200 {array} model.TimelineItem
// @Failure 400 {object} dto.MessageResponse
// @Router /timeline [get]
func (h *Controller) getCommonTimeline(c *gin.Context) {
	timeline, err := h.services.Timeline.GetCommon()
	commonJSONResponse(c, timeline, err)
}

// @Summary Get personalized timeline
// @Security Auth
// @Tags timeline
// @Description Получение персонализированного таймлайна. Состоит из событий тех пользователей, на которых подписан текущий
// @Accept json
// @Produce json
// @Success 200 {array} model.TimelineItem
// @Failure 400,403 {object} dto.MessageResponse
// @Router /timeline/personalized [get]
func (h *Controller) getPersonalizedTimeline(c *gin.Context) {
	currentUser, _ := getUserId(c)
	timeline, err := h.services.Timeline.GetForUser(currentUser)
	commonJSONResponse(c, timeline, err)
}

// @Summary Get own timeline
// @Security Auth
// @Tags timeline
// @Description Получение личного таймлайна. Состоит из личных событий текущего пользователя (включая анонимные)
// @Accept json
// @Produce json
// @Success 200 {array} model.TimelineItem
// @Failure 400,403 {object} dto.MessageResponse
// @Router /timeline/own [get]
func (h *Controller) getOwnTimeline(c *gin.Context) {
	currentUser, _ := getUserId(c)
	timeline, err := h.services.Timeline.GetUserOwn(currentUser)
	commonJSONResponse(c, timeline, err)
}
