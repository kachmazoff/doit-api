package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Controller) userIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
	}

	c.Set(userCtx, id)
}

func (h *Controller) optionalUserIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err == nil {
		c.Set(userCtx, id)
	}
}

func (h *Controller) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func getUserId(c *gin.Context) (string, error) {
	idFromCtx, ok := c.Get(userCtx)
	if !ok {
		return "", errors.New("userCtx not found")
	}

	id, ok := idFromCtx.(string)
	if !ok {
		return "", errors.New("userCtx is of invalid type")
	}

	return id, nil
}
