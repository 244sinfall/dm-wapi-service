package controllers

import (
	services "darkmoon-wapi-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateLottery(c *gin.Context) {
	var lotteryCreator services.LotteryOptions
	if err := c.BindJSON(&lotteryCreator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if lotteryCreator.Rate < 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lottery could be created for rate 7 and more"})
		return
	}
	if lotteryCreator.ParticipantsCount < 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lottery could be created for more than 10 participants."})
		return
	}
	respond := lotteryCreator.GenerateLottery()
	c.JSON(http.StatusOK, respond)
}
