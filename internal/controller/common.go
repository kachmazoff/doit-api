package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/dto"
)

func createMessage(message string) dto.MessageResponse {
	return dto.MessageResponse{Message: message}
}

func commonJSONResponse(c *gin.Context, data interface{}, err error) {
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, data)
}

func handleCreation(c *gin.Context, id string, err error) {
	response := dto.IdResponse{Id: id}
	commonJSONResponse(c, response, err)
}

func handleUserFindError(c *gin.Context, err error) bool {
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, createMessage("This user does not exist"))
			return true
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return true
	}
	return false
}
