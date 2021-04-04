package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/model"
)

func (h *Controller) initTimelineRoutes(api *gin.RouterGroup) {
	timeline := api.Group("/timeline")
	{
		timeline.GET("/", h.optionalUserIdentity, h.getTimeline)
		timeline.GET("/common", h.getCommonTimeline)
		timeline.GET("/own", h.userIdentity, h.getOwnTimeline)
		timeline.GET("/personalized", h.userIdentity, h.getPersonalizedTimeline)
	}
}

func (h *Controller) getCommonTimeline(c *gin.Context) {
	timeline, err := h.services.Timeline.GetCommon()
	commonJSONResponse(c, timeline, err)
}

// @Summary Get timeline with filters
// @Tags timeline
// @Description Получение таймлайна по фильтрам
// @Accept json
// @Produce json
// @Param userId query string false "Id пользователя"
// @Param type query string false "Тип запрашиваемого таймлайна (subs, common)"
// @Param participantId query string false "Id дневника/участника"
// @Param challengeId query string false "Id челленджа"
// @Param eventTypes query []string false "Массив типов событий ('CREATE_CHALLENGE', 'ACCEPT_CHALLENGE', 'ADD_NOTE', 'ADD_SUGGESTION')"
// @Param order query string false "Порядок сортировки (ASC, DESC)"
// @Param lastIndex query int false "Индекс последней полученной записи"
// @Param limit query int false "Максимальное количество записей"
// @Success 200 {array} model.TimelineItem
// @Failure 400 {object} dto.MessageResponse
// @Router /timeline [get]
func (h *Controller) getTimeline(c *gin.Context) {
	currentUser, err := getUserId(c)
	limit, _ := strconv.Atoi(c.Query("limit"))

	if limit > 100 {
		limit = 100
	}
	if limit <= 0 {
		limit = 20
	}

	lastIndex, err := strconv.Atoi(c.Query("lastIndex"))
	if err != nil {
		lastIndex = -1
	}

	orderType := c.Query("order")
	if orderType != "asc" && orderType != "desc" {
		orderType = "desc"
	}
	orderType = strings.ToUpper(orderType)

	userId := c.Query("userId")

	timelineType := c.Query("type")
	if timelineType == "" {
		timelineType = "common"
	} else if timelineType == "subs" {
		if currentUser == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}
		if userId != "" && userId != currentUser {
			c.AbortWithStatusJSON(http.StatusForbidden, err)
			return
		} else {
			userId = currentUser
		}
	}

	filters := model.TimelineFilters{
		RequestAuthor: currentUser,
		UserId:        userId,
		ParticipantId: c.Query("participantId"),
		ChallengeId:   c.Query("challengeId"),
		Limit:         limit,
		LastIndex:     lastIndex,
		EventTypes:    c.QueryArray("eventTypes"),
		TimelineType:  timelineType,
		OrderType:     orderType,
	}

	timeline, err := h.services.Timeline.GetWithFilters(filters)
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
