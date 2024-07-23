package controllers

import (
	services "darkmoon-wapi-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CleanParticipantsText(c *gin.Context) {
	var raw services.ParticipantsRequest
	err := c.BindJSON(&raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	response := raw.CleanRawText()
	c.JSON(http.StatusOK, response)
}
