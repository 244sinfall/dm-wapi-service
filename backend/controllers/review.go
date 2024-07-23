package controllers

import (
	services "darkmoon-wapi-service/services"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Review = services.Review

func ReviewGenerate(c *gin.Context) {
	var reviewObject Review

	if err := c.BindJSON(&reviewObject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(reviewObject.Rates) == 0 || reviewObject.CharName == "" ||
		reviewObject.ReviewerDiscord == "" || reviewObject.ReviewerProfile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not all required fields provided"})
	}

	c.JSON(http.StatusOK, reviewObject.GetReviewResponse())
}
