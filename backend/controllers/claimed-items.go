package controllers

import (
	"context"
	"encoding/json"

	"darkmoon-wapi-service/permissions"
	rabbitmqclient "darkmoon-wapi-service/rabbitmq_client"
	rabbitmqmessages "darkmoon-wapi-service/rabbitmq_messages"
	services "darkmoon-wapi-service/services"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type ClaimedItem = services.ClaimedItem

type ClaimedItemsList = services.ClaimedItemsList

func AddClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	user, err := services.Authenticate(c.Request.Header.Get("Authorization"), a, f, ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user.Permission < permissions.GmPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	claimedItem := new(ClaimedItem)
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(claimedItem)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	list := services.GetClaimedItems()
	err = list.Add(*claimedItem, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(200)
		baseMessage := rabbitmqmessages.GetAuthenticatedBaseMessage(user, "Claimed-Item", "Added")
		rabbitmqclient.SendLogMessage(rabbitmqmessages.GetClaimedItemRabbitMQMessage(baseMessage, claimedItem, nil), user.AuthenticatedUser, nil)
		return
	}
}

func ApproveClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	user, err := services.Authenticate(c.Request.Header.Get("Authorization"), a, f, ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user.Permission < permissions.AdminPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	id := c.Param("id")
	fbUser, err := a.GetUser(ctx, user.IntegrationUserId)
	if err != nil {
		c.JSON(500, gin.H{"error": "User not found"})
		return
	}
	list := services.GetClaimedItems()
	err = list.Approve(id, fbUser.DisplayName, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(200)
		return
	}
}

func UpdateClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	user, err := services.Authenticate(c.Request.Header.Get("Authorization"), a, f, ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user.Permission < permissions.GmPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	fbUser, _ := a.GetUser(ctx, user.IntegrationUserId)
	claimedItemMock := new(ClaimedItem)
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(claimedItemMock)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	list := services.GetClaimedItems()
	newItem, oldItem, err := list.Update(id, user.Permission >= permissions.AdminPermission, *claimedItemMock, fbUser.DisplayName, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		baseMessage := rabbitmqmessages.GetAuthenticatedBaseMessage(user, "Claimed-Item", "Updated")
		rabbitmqclient.SendLogMessage(rabbitmqmessages.GetClaimedItemRabbitMQMessage(baseMessage, newItem, oldItem), user.AuthenticatedUser, nil)
		c.Status(200)
		return
	}
}

func DeleteClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	user, err := services.Authenticate(c.Request.Header.Get("Authorization"), a, f, ctx)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if user.Permission < permissions.AdminPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	id := c.Param("id")
	list := services.GetClaimedItems()
	item, err := list.Delete(id, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	baseMessage := rabbitmqmessages.GetAuthenticatedBaseMessage(user, "Claimed-Item", "Deleted")
	rabbitmqclient.SendLogMessage(rabbitmqmessages.GetClaimedItemRabbitMQMessage(baseMessage, item, nil), user.AuthenticatedUser, nil)
	c.Status(200)
}

func ReceiveClaimedItems(c *gin.Context) {
	c.JSON(200, gin.H{"result": services.GetClaimedItems()})
}
