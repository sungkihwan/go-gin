package handlers

import (
	"go-gin-postgre/domain"
	"go-gin-postgre/domain/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userHandler struct {
	usecase usecases.UserUsecase
}

func NewUserHandler(u usecases.UserUsecase) UserHandler {
	return &userHandler{u}
}

func (h *userHandler) GetUsers(c *gin.Context) {
	users, err := h.usecase.GetUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}

func (h *userHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")

	// 문자열 id를 uint로 변환
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.usecase.GetUserByID(uint(id))
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}
	err := h.usecase.CreateUser(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}
	c.JSON(201, user)
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}
	err := h.usecase.UpdateUser(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error updating user"})
		return
	}
	c.JSON(200, user)
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID format"})
		return
	}
	err = h.usecase.DeleteUser(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting user"})
		return
	}
	c.JSON(204, nil)
}
