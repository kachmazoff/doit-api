package controller

import "github.com/gin-gonic/gin"

func (h *Controller) initTimelineRoutes(api *gin.RouterGroup) {
	timeline := api.Group("/timeline")
	{
		timeline.GET("/", h.getAll)
		timeline.GET("/own", h.userIdentity, h.getOwnTimeline)
		timeline.GET("/personalized", h.userIdentity, h.getPersonalizedTimeline)
	}
}

func (h *Controller) getAll(c *gin.Context) {
	timeline, err := h.services.Timeline.GetCommon()
	commonJSONResponse(c, timeline, err)
}

func (h *Controller) getPersonalizedTimeline(c *gin.Context) {
	currentUser, _ := getUserId(c)
	timeline, err := h.services.Timeline.GetForUser(currentUser)
	commonJSONResponse(c, timeline, err)
}

func (h *Controller) getOwnTimeline(c *gin.Context) {
	currentUser, _ := getUserId(c)
	timeline, err := h.services.Timeline.GetUserOwn(currentUser)
	commonJSONResponse(c, timeline, err)
}
