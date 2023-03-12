package main

import (
	"context"
	"darkmoonWebApi/arbiters"
	"darkmoonWebApi/charsheet"
	claimed_items "darkmoonWebApi/claimed-items"
	"darkmoonWebApi/economics"
	"darkmoonWebApi/events"
	"darkmoonWebApi/gob"
	"darkmoonWebApi/other"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
	"os"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, DELETE, GET, PUT, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	var opt = option.WithCredentialsFile(os.Getenv("DM_API_FIREBASE_CREDENTIALS_FILE"))
	var ctx = context.Background()
	var app, err = firebase.NewApp(ctx, nil, opt)
	firestore, err := app.Firestore(ctx)
	auth, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error initializing firebase %v\n", err)
	}
	go economics.ChecksScheduler(false)
	claimed_items.GetClaimedItemsFromDatabase(firestore, ctx)
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
		economics.ReceiveChecks(c, auth, firestore, ctx)
	})
	router.GET("/gobs", func(c *gin.Context) {
		gob.ReceiveGobs(c, auth, firestore, ctx)
	})
	router.GET("/claimed_items/get_items", claimed_items.ReceiveClaimedItems)
	router.DELETE("/claimed_items/delete/:id", func(c *gin.Context) {
		claimed_items.DeleteClaimedItem(c, auth, firestore, ctx)
	})
	router.PUT("/claimed_items/update/:id", func(c *gin.Context) {
		claimed_items.UpdateClaimedItem(c, auth, firestore, ctx)
	})
	router.PATCH("/claimed_items/approve/:id", func(c *gin.Context) {
		claimed_items.ApproveClaimedItem(c, auth, firestore, ctx)
	})
	router.POST("/claimed_items/create", func(c *gin.Context) {
		claimed_items.AddClaimedItem(c, auth, firestore, ctx)
	})
	err = router.RunTLS("185.193.143.35:8443", "cert.pem", "privkey.pem")
	//err = router.Run("127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
	}

}
