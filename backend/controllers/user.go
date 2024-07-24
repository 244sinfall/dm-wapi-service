package controllers

import (
	"context"
	permissions "darkmoon-wapi-service/permissions"
	services "darkmoon-wapi-service/services"
	"encoding/json"
	"errors"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func ResetUserPassword(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	user, err := services.Authenticate(c.Request.Header.Get("Authorization"), a, f, ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user.Permission < permissions.AdminPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	body := new(struct {
		Email string `json:"email"`
	})
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err = a.PasswordResetLink(ctx, body.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

type connectBody struct {
	Code string `json:"code"`
}

func ConnectToAuthService(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	fbaccess := c.Request.Header.Get("Authorization")
	decoder := json.NewDecoder(c.Request.Body)
	var bodyJson = new(connectBody)
	err := decoder.Decode(&bodyJson)
	if err != nil {
		c.JSON(400, gin.H{"error": "Error decoding body: " + err.Error()})
		return
	}
	user, err := services.ConnectToDarkmoon(fbaccess, bodyJson.Code, a, f, ctx)
	if err != nil {
		c.JSON(503, gin.H{"error": "Error connecting to Darkmoon: " + err.Error()})
		return
	}
	c.JSON(200, user)
}

func GetMe(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	fbaccess := c.Request.Header.Get("Authorization")
	user, err := services.Authenticate(fbaccess, a, f, ctx)
	if err != nil {
		if errors.Is(err, &services.NotConnectedError{}) || errors.Is(err, &services.RevokedError{}) {
			c.JSON(404, gin.H{"error": err.Error()})
		}
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}
