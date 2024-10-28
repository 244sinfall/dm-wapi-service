package main

import (
	"darkmoon-wapi-service/arbiter"
	"darkmoon-wapi-service/auth"
	"darkmoon-wapi-service/checks"
	claimeditems "darkmoon-wapi-service/claimed-items"
	"darkmoon-wapi-service/gob"
	logcleaner "darkmoon-wapi-service/log-cleaner"
	"darkmoon-wapi-service/participants"
	"darkmoon-wapi-service/review"
	"fmt"
	"log"
	"os"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
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
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           os.Getenv("SENTRY_DSN"),
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{}))
	router.Use(CORSMiddleware())
	router.POST("/generate_charsheet_review", review.GenerateReview)
	router.POST("/events/clean_participants_text", participants.CleanParticipantsText)
	router.POST("/events/create_lottery", review.CreateLottery)
	router.POST("/arbiters/rewards_work", arbiter.GenerateArbiterCommands)
	router.POST("/clean_log", logcleaner.CleanLog)
	router.GET("/economics/get_checks", checks.ReceiveChecks)
	router.GET("/gobs", gob.ReceiveGobs)
	router.GET("/claimed_items/get_items", claimeditems.ReceiveClaimedItems)
	router.DELETE("/claimed_items/delete/:id", claimeditems.DeleteClaimedItem)
	router.PUT("/claimed_items/update/:id", claimeditems.UpdateClaimedItem)
	router.PATCH("/claimed_items/approve/:id", claimeditems.ApproveClaimedItem)
	router.POST("/claimed_items/create", claimeditems.AddClaimedItem)
	router.POST("/users/reset", auth.ResetUserPassword)
	router.POST("v2/users/connect", auth.ConnectToAuthService)
	router.GET("v2/users/me", auth.GetMe)
	err := router.Run("0.0.0.0:" + os.Getenv("BACKEND_PORT"))
	if err != nil {
		log.Fatalf("Error on listening: %v", err)
	}

}
