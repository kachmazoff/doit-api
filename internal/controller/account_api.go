package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kachmazoff/doit-api/internal/dto"
	"github.com/kachmazoff/doit-api/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (h *Controller) initAccountRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/registration", h.registerUser)
		auth.POST("/login", h.getToken)
		auth.POST("/activate", h.activateAccount)
	}
}

// @Summary Registration
// @Tags auth
// @Description Создание нового пользователя
// @Accept json
// @Produce json
// @Param input body dto.Registration true "Данные пользователя"
// @Success 200 {object} dto.IdResponse
// @Router /auth/registration [post]
func (h *Controller) registerUser(c *gin.Context) {
	var input dto.Registration
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	userModel := model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	id, err := h.services.Users.Create(userModel)
	handleCreation(c, id, err)
}

// @Summary Account activation
// @Tags auth
// @Description Активация нового аккаунта (подтверждение электронной почты)
// @Accept json
// @Produce json
// @Success 200 {object} dto.MessageResponse
// @Router /auth/activate [post]
func (h *Controller) activateAccount(c *gin.Context) {
	userId := c.Query("id")

	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Bad request"))
		return
	}

	err := h.services.Users.ConfirmAccount(userId)
	commonJSONResponse(c, createMessage("Account activated successfully"), err)
}

// @Summary Login
// @Tags auth
// @Description Получение jwt-токена для дальнейшей работы с сервисом
// @Accept json
// @Produce json
// @Param input body dto.Login true "Данные пользователя"
// @Success 200 {object} dto.TokenResponse
// @Failure 400,404 {object} dto.MessageResponse
// @Router /auth/login [post]
func (h *Controller) getToken(c *gin.Context) {
	var input dto.Login
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user, err := h.services.Users.GetByEmail(input.Email)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Пользователя с данным email не существует"))
		return
	}

	password := []byte(input.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 6)

	if string(hashedPassword) != user.Password {
		c.AbortWithStatusJSON(http.StatusBadRequest, createMessage("Неверный пароль"))
		return
	}

	token, err := h.tokenManager.NewJWT(user.Id, time.Hour)

	commonJSONResponse(c, dto.TokenResponse{Token: token}, err)
}
