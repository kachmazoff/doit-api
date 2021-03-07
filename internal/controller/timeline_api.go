package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Controller) initTimelineRoutes(api *gin.RouterGroup) {
	timeline := api.Group("/timeline")
	{
		timeline.GET("/", h.getAll)
	}
}

func (h *Controller) getAll(c *gin.Context) {
	timeline, err := h.services.Timeline.GetCommon()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, timeline)
}
