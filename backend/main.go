package main

import (
	"context"
	controllers "darkmoon-wapi-service/controllers"
	services "darkmoon-wapi-service/services"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
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
	if err != nil {
		log.Printf("error initializing firebase %v\n", err)
	}
	auth, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error initializing firebase %v\n", err)
	}
	go services.ChecksScheduler(false)
	services.GetClaimedItemsFromDatabase(firestore, ctx)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())
	// *Review Generator*
	router.POST("/generate_charsheet_review", func(c *gin.Context) {
		controllers.ReviewGenerate(c)
	})
	// *Events* Clean Text
	router.POST("/events/clean_participants_text", func(c *gin.Context) {
		controllers.CleanParticipantsText(c)
	})
	// *Events* Create Lottery
	router.POST("/events/create_lottery", func(c *gin.Context) {
		controllers.CreateLottery(c)
	})
	// *Arbiters* Base
	router.POST("/arbiters/rewards_work", func(c *gin.Context) {
		controllers.ArbiterCalculation(c)
	})
	// *Log Cleaner*
	router.POST("/clean_log", func(c *gin.Context) {
		controllers.LogClean(c)
	})
	// *Economics* Get Checks
	router.GET("/economics/get_checks", func(c *gin.Context) {
		controllers.ReceiveChecks(c, auth, firestore, ctx)
	})
	router.GET("/gobs", func(c *gin.Context) {
		controllers.ReceiveGobs(c, auth, firestore, ctx)
	})
	router.GET("/claimed_items/get_items", func(c *gin.Context) {
		controllers.ReceiveClaimedItems(c)
	})
	router.DELETE("/claimed_items/delete/:id", func(c *gin.Context) {
		controllers.DeleteClaimedItem(c, auth, firestore, ctx)
	})
	router.PUT("/claimed_items/update/:id", func(c *gin.Context) {
		controllers.UpdateClaimedItem(c, auth, firestore, ctx)
	})
	router.PATCH("/claimed_items/approve/:id", func(c *gin.Context) {
		controllers.ApproveClaimedItem(c, auth, firestore, ctx)
	})
	router.POST("/claimed_items/create", func(c *gin.Context) {
		controllers.AddClaimedItem(c, auth, firestore, ctx)
	})
	router.POST("/users/reset", func(c *gin.Context) {
		controllers.ResetUserPassword(c, auth, firestore, ctx)
	})
	err = router.Run("0.0.0.0:80")
	if err != nil {
		fmt.Println(err)
	}

}
