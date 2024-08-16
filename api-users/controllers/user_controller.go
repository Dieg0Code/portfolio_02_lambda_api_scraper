package controllers

import "github.com/gin-gonic/gin"

type UserController interface {
	RegisterUser(c *gin.Context)
	GetAllUsers(c *gin.Context)
	GetUserByID(c *gin.Context)
	LogInUser(c *gin.Context)
}
