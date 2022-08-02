package main

import (
	"darkmoonWebApi/charsheet"
	"darkmoonWebApi/events"
	"darkmoonWebApi/other"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())
	// *Review Generator*
	router.POST("/generate_charsheet_review", charsheet.GenerateReview)
	// *Events* Clean Text
	router.POST("/events/clean_participants_text", events.CleanParticipantsText)
	// *Events* Create Lottery
	router.POST("/events/create_lottery", events.CreateLottery)
	// *Log Cleaner*
	router.POST("/clean_log", other.CleanLog)

	//err := router.RunTLS("dm.rolevik.site:8443", "cert.pem", "privkey.pem")
	err := router.Run("127.0.0.1:8000")
	if err != nil {
		return
	}

}
