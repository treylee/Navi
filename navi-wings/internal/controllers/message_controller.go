package controllers

import (
	"navi-wings/config"
	"navi-wings/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {
	var messages []models.Message
	config.DB.Find(&messages)
	c.JSON(http.StatusOK, messages)
}

func CreateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&message)
	c.JSON(http.StatusOK, message)
}
