package participants

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CleanParticipantsText(c *gin.Context) {
	var raw participantsRequest
	err := c.BindJSON(&raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	response := raw.cleanRawText()
	c.JSON(http.StatusOK, response)
}
