package user

import (
	"context"
	"net/http"
	"strconv"
	"time"

	userModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/user"
	userService "github.com/Ayocodes24/GO-Eats/pkg/service/user"
	"github.com/gin-gonic/gin"
)

func (s *UserHandler) addUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var u userModel.User
	if err := c.BindJSON(&u); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := s.validate.Struct(u); err != nil {
		validationError := userModel.UserValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": validationError})
		return
	}

	_, err := s.service.Add(ctx, &u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (s *UserHandler) deleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	idParam := c.Param("id")
	userID, _ := strconv.ParseInt(idParam, 10, 64)

	_, err := s.service.Delete(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *UserHandler) loginUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var creds userModel.LoginUser
	if err := c.BindJSON(&creds); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	login := userService.ValidateAccount(
		s.service.Login,
		s.service.UserExist,
		s.service.ValidatePassword,
	)
	token, err := login(ctx, &userModel.LoginUser{
		Email:    creds.Email,
		Password: creds.Password,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
