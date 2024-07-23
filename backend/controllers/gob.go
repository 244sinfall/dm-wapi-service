package controllers

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	services "darkmoon-wapi-service/services"
)


func ReceiveGobs(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, _ := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	permInfo, _ := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if permission < gmPermission {
		c.JSON(400, gin.H{"error": "Not enough permissions"})
		return
	}
	c.JSON(200, gin.H{"result": services.GetGobs()})
}
