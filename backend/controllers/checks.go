package controllers

import (
	"context"
	"strconv"
	"time"

	"darkmoon-wapi-service/permissions"
	services "darkmoon-wapi-service/services"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func ReceiveChecks(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	search := c.Query("search")
	category := c.Query("category")
	status := c.Query("status")
	sortMethod := c.Query("sortBy")           //
	sortDirection := c.Query("sortDirection") //
	force := c.Query("force")
	CachedChecks := services.GetCachedChecks()
	if CachedChecks.Updating {
		c.JSON(500, gin.H{"error": "Checks are currently unavailable due to cache update"})
		return
	}
	if len(CachedChecks.Checks) == 0 {
		services.ChecksScheduler(true)
	}
	if force != "" {
		user, err := services.Authenticate(c.Request.Header.Get("Authorization"), a, f, ctx)
		if err != nil {
			c.JSON(403, gin.H{"error": "You don't have permission"})
			return
		}
		if user.Permission < permissions.GmPermission {
			c.JSON(403, gin.H{"error": "Not enough permissions"})
			return
		}
		if time.Since(CachedChecks.UpdatedAt) < 5*time.Minute {
			c.JSON(400, gin.H{"error": "Force update is available if cached checks are older than 5 minutes", "updatedAt": CachedChecks.UpdatedAt})
			return
		} else {
			err := services.ParseAndDeployNewChecks()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}
	sentChecks := CachedChecks.Checks
	if category != "" {
		sentChecks = services.FilterChecksCategory(sentChecks, category)
	}
	if status != "" {
		sentChecks = services.FilterChecksStatus(sentChecks, status)
	}
	if search != "" {
		sentChecks = services.FilterChecks(sentChecks, search)
	}
	filteredCount := len(sentChecks)
	if skip != 0 {
		if len(sentChecks)-1 > skip {
			sentChecks = sentChecks[skip:]
		} else {
			sentChecks = []services.CheckResponseCheck{}
		}
	}
	if limit != 0 {
		if len(sentChecks) > limit {
			sentChecks = sentChecks[:limit]
		}
	}
	var sortDir bool
	if sortDirection == "ascending" {
		sortDir = true
	}
	services.SortChecks(sentChecks, sortMethod, sortDir)
	c.JSON(200, gin.H{"checks": sentChecks, "count": len(CachedChecks.Checks),
		"filteredCount": filteredCount, "updatedAt": CachedChecks.UpdatedAt, "types": CachedChecks.Types})
}
