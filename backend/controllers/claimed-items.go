package controllers

import (
	"context"
	"encoding/json"

	services "darkmoon-wapi-service/services"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type ClaimedItem = services.ClaimedItem

type ClaimedItemsList = services.ClaimedItemsList

func AddClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	permInfo, err := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if permission < reviewerPermission {
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
		return
	}
}

func ApproveClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	permInfo, err := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if permission < adminPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	id := c.Param("id")
	user, err := a.GetUser(ctx, token.UID)
	if err != nil {
		c.JSON(500, gin.H{"error": "User not found"})
		return
	}
	list := services.GetClaimedItems()
	err = list.Approve(id, user.DisplayName, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(200)
		return
	}
}

func UpdateClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	permInfo, err := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if permission < reviewerPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	user, _ := a.GetUser(ctx, token.UID)
	claimedItemMock := new(ClaimedItem)
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(claimedItemMock)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	list := services.GetClaimedItems()
	err = list.Update(id, permission == adminPermission, *claimedItemMock, user.DisplayName, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(200)
		return
	}
}

func DeleteClaimedItem(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	permInfo, err := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if permission < adminPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	id := c.Param("id")
	list := services.GetClaimedItems()
	err = list.Delete(id, f, ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

func ReceiveClaimedItems(c *gin.Context) {
	c.JSON(200, gin.H{"result": services.GetClaimedItems()})
}
