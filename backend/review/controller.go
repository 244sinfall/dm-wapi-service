package review

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenerateReview(c *gin.Context) {
	var reviewObject review

	if err := c.BindJSON(&reviewObject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(reviewObject.Rates) == 0 || reviewObject.CharName == "" ||
		reviewObject.ReviewerDiscord == "" || reviewObject.ReviewerProfile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not all required fields provided"})
	}

	c.JSON(http.StatusOK, reviewObject.getReviewResponse())
}

func CreateLottery(c *gin.Context) {
	var lotteryCreator lotteryOptions
	if err := c.BindJSON(&lotteryCreator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if lotteryCreator.Rate < minRate {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough rate"})
		return
	}
	if lotteryCreator.ParticipantsCount < minParticipants {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough participants."})
		return
	}
	respond := lotteryCreator.generateLottery()
	c.JSON(http.StatusOK, respond)
}
