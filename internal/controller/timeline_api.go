package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

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

	timelineResponse(c, timeline, err)
}

func (h *Controller) getPersonalizedTimeline(c *gin.Context) {
	currentUser, _ := getUserId(c)
	timeline, err := h.services.Timeline.GetForUser(currentUser)
	timelineResponse(c, timeline, err)
}

func (h *Controller) getOwnTimeline(c *gin.Context) {
	currentUser, _ := getUserId(c)
	timeline, err := h.services.Timeline.GetUserOwn(currentUser)
	timelineResponse(c, timeline, err)
}

func timelineResponse(c *gin.Context, timeline []model.TimelineItem, err error) {
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, timeline)
}
