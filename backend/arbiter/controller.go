package arbiter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenerateArbiterCommands(c *gin.Context) {
	var request arbiterCommandsRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch request.Mode {
	case "givexp", "givegold", "takexp":
		response := request.generateResponse()
		c.JSON(http.StatusOK, response)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown arbiter work mode"})
		return
	}
}