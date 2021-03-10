package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createMessage(message string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
	}
}

func handleCreation(c *gin.Context, id string, err error) {
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, createMessage(err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
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
