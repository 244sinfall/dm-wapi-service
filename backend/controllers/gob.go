package controllers

import (
	"context"
	"darkmoon-wapi-service/permissions"
	services "darkmoon-wapi-service/services"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)


func ReceiveGobs(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	user, err := services.Authenticate(c.Request.Header.Get("Authorization"), a, f, ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user.Permission < permissions.GmPermission {
		c.JSON(400, gin.H{"error": "Not enough permissions"})
		return
	}
	c.JSON(200, gin.H{"result": services.GetGobs()})
}
