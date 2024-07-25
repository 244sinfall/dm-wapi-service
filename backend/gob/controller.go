package gob

import (
	"darkmoon-wapi-service/auth"

	"github.com/gin-gonic/gin"
)

func ReceiveGobs(c *gin.Context) {
	user, err := auth.Authenticate(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !user.IsGM() {
		c.JSON(400, gin.H{"error": "Not enough permissions"})
		return
	}
	c.JSON(200, gin.H{"result": getGobs()})
}
