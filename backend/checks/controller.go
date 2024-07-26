package checks

import (
	"time"

	"github.com/gin-gonic/gin"
)

func ReceiveChecks(c *gin.Context) {
	filter := new(checkRequestFilter)
	err := c.BindQuery(&filter)
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad query: " + err.Error()})
		return
	}
	checks, err := getChecks(filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed at getting checks: " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"checks": checks.Result, "count": checks.Count,
		"filteredCount": checks.Count, "updatedAt": time.Now(), "types": checks.Types})
}
