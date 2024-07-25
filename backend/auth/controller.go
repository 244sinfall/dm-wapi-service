package auth

import (
	"darkmoon-wapi-service/globals"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

func ResetUserPassword(c *gin.Context) {
	auth := globals.GetAuth()
	authHeader := c.Request.Header.Get("Authorization")
	token, err := auth.VerifyIDToken(globals.GetGlobalContext(), authHeader)
	if err != nil {
		c.JSON(401, gin.H{"error": "Error auth: " + err.Error()})
		return
	}
	user, err := auth.GetUser(globals.GetGlobalContext(), token.UID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err = auth.PasswordResetLink(globals.GetGlobalContext(), user.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}


func ConnectToAuthService(c *gin.Context) {
	auth := globals.GetAuth()
	fbaccess := c.Request.Header.Get("Authorization")
	token, err := auth.VerifyIDToken(globals.GetGlobalContext(), fbaccess)
	if err != nil {
		c.JSON(401, gin.H{"error": "Error auth: " + err.Error()})
		return
	}
	decoder := json.NewDecoder(c.Request.Body)
	var bodyJson = new(connectRequestBody)
	err = decoder.Decode(&bodyJson)
	if err != nil {
		c.JSON(400, gin.H{"error": "Error decoding body: " + err.Error()})
		return
	}
	user, err := connectToDarkmoon(token.UID, bodyJson.Code)
	if err != nil {
		c.JSON(503, gin.H{"error": "Error connecting to Darkmoon: " + err.Error()})
		return
	}
	c.JSON(200, user)
}

func GetMe(c *gin.Context) {
	fbaccess := c.Request.Header.Get("Authorization")
	user, err := Authenticate(fbaccess)
	if err != nil {
		if errors.Is(err, &notConnectedError{}) || errors.Is(err, &revokedError{}) {
			c.JSON(404, gin.H{"error": err.Error()})
		}
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, user)
}
