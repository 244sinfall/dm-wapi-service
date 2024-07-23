package controllers

import (
	services "darkmoon-wapi-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ArbiterCalculation(c *gin.Context) {
	var request services.ArbiterCalculationRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	switch request.Mode {
	case "givexp", "givegold", "takexp":
		response := request.GenerateResponse()
		c.JSON(http.StatusOK, response)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown arbiter work mode"})
		return
	}
}
