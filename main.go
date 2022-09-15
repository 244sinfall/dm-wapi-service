package main

import (
	"context"
	"darkmoonWebApi/arbiters"
	"darkmoonWebApi/charsheet"
	"darkmoonWebApi/economics"
	"darkmoonWebApi/events"
	"darkmoonWebApi/other"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
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
	var opt = option.WithCredentialsFile("darkmoon-web-api-2-firebase-adminsdk-7x7ol-b33aadf8c6.json")
	var ctx = context.Background()
	var app, err = firebase.NewApp(ctx, nil, opt)
	firestore, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("error initializing firebase %v\n", err)
	}
	go economics.ChecksScheduler(firestore, ctx, false)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())
	// *Review Generator*
	router.POST("/generate_charsheet_review", charsheet.GenerateReview)
	// *Events* Clean Text
	router.POST("/events/clean_participants_text", events.CleanParticipantsText)
	// *Events* Create Lottery
	router.POST("/events/create_lottery", events.CreateLottery)
	// *Arbiters* Base
	router.POST("/arbiters/rewards_work", arbiters.ArbiterWork)
	// *Log Cleaner*
	router.POST("/clean_log", other.CleanLog)
	// *Economics* Get Checks
	router.GET("/economics/get_checks", func(c *gin.Context) {
		economics.ReceiveChecks(c, firestore, ctx)
	})
	//err := router.RunTLS("dm.rolevik.site:8443", "cert.pem", "privkey.pem")
	err = router.Run("127.0.0.1:8000")
	if err != nil {
		return
	}

}
