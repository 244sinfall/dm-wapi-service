package economics

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type APIResponseItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type APIResponseCheck struct {
	Id       int               `json:"id"`
	Date     string            `json:"date"`
	Sender   string            `json:"sender"`   // owner
	Receiver string            `json:"receiver"` // checktype
	Subject  string            `json:"subject"`  // name
	Body     string            `json:"body"`     // description
	Money    int               `json:"money"`
	GmName   string            `json:"gmName"`
	Status   string            `json:"status"`
	Items    []APIResponseItem `json:"items"`
}

const gmPermission = 1

const defaultCheckCount = 13000

func ParseChecksFromDarkmoon() ([]APIResponseCheck, error) {
	newChecks := make([]APIResponseCheck, 0, defaultCheckCount)
	response, err := http.Get(os.Getenv("DM_API_CHECKS_ADDRESS"))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&newChecks)
	if len(newChecks) == 0 {
		return nil, errors.New("nothing parsed")
	}
	if err != nil {
		return []APIResponseCheck{}, err
	}
	return newChecks, nil
}
func filterChecksCategory(c []APIResponseCheck, category string) []APIResponseCheck {
	newChecks := make([]APIResponseCheck, 0, len(c)/5)
	for _, check := range c {
		if check.Receiver == category {
			newChecks = append(newChecks, check)
		}
	}
	return newChecks
}

func filterChecksStatus(c []APIResponseCheck, status string) []APIResponseCheck {
	var statusName string
	switch status {
	case "open":
		statusName = "Ожидает"
	case "closed":
		statusName = "Закрыт"
	case "rejected":
		statusName = "Отказан"
	default:
		return c
	}
	newChecks := make([]APIResponseCheck, 0, len(c)/5)
	for _, check := range c {
		if check.Status == statusName {
			newChecks = append(newChecks, check)
		}
	}
	return newChecks
}

func sortChecks(c []APIResponseCheck, sortBy string, ascending bool) {
	switch sortBy {
	case "money":
		sort.Slice(c, func(i, j int) bool {
			if ascending {
				return c[i].Money < c[j].Money
			}
			return c[i].Money > c[j].Money
		})
	case "date":
		sort.Slice(c, func(i, j int) bool {
			iTime, _ := time.Parse("02.01.2006 15:04", c[i].Date)
			jTime, _ := time.Parse("02.01.2006 15:04", c[j].Date)
			if ascending {
				return iTime.Unix() < jTime.Unix()
			}
			return iTime.Unix() > jTime.Unix()
		})
	default:
		sort.Slice(c, func(i, j int) bool {
			if ascending {
				return c[i].Id < c[j].Id
			}
			return c[i].Id > c[j].Id
		})
	}
}

func findTextMatch(str string, phrase string) bool {
	lowerStr := strings.ToLower(str)
	if strings.Contains(lowerStr, phrase) {
		return true
	}
	return false
}

func filterChecks(c []APIResponseCheck, phrase string) []APIResponseCheck {
	if len(c) == 0 {
		return c
	}
	newChecks := make([]APIResponseCheck, 0, len(c)/10)
	lowerPhrase := strings.ToLower(phrase)
	for _, check := range c {
		if strings.Contains(strconv.Itoa(check.Id), phrase) || findTextMatch(check.GmName, lowerPhrase) ||
			findTextMatch(check.Subject, lowerPhrase) || findTextMatch(check.Sender, lowerPhrase) ||
			findTextMatch(check.Body, lowerPhrase) || findTextMatch(check.Receiver, lowerPhrase) {
			newChecks = append(newChecks, check)
		}
	}
	return newChecks
}

func ReceiveChecks(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	search := c.Query("search")
	category := c.Query("category")
	status := c.Query("status")
	sortMethod := c.Query("sortBy")           //
	sortDirection := c.Query("sortDirection") //
	force := c.Query("force")
	if CachedChecks.updating {
		c.JSON(500, gin.H{"error": "Checks are currently unavailable due to cache update"})
		return
	}
	if len(CachedChecks.checks) == 0 {
		ChecksScheduler(true)
	}
	if force != "" {
		token, _ := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
		permInfo, _ := f.Doc("permissions/" + token.UID).Get(ctx)
		permission := permInfo.Data()["permission"].(int64)
		if permission < gmPermission {
			c.JSON(400, gin.H{"error": "Not enough permissions"})
			return
		}
		if time.Now().Sub(CachedChecks.updatedAt) < 5*time.Minute {
			c.JSON(400, gin.H{"error": "Force update is available if cached checks are older than 5 minutes", "updatedAt": CachedChecks.updatedAt})
			return
		} else {
			err := ParseAndDeployNewChecks()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}
	sentChecks := CachedChecks.checks
	if category != "" {
		sentChecks = filterChecksCategory(sentChecks, category)
	}
	if status != "" {
		sentChecks = filterChecksStatus(sentChecks, status)
	}
	if search != "" {
		sentChecks = filterChecks(sentChecks, search)
	}
	filteredCount := len(sentChecks)
	if skip != 0 {
		if len(sentChecks)-1 > skip {
			sentChecks = sentChecks[skip:]
		} else {
			sentChecks = []APIResponseCheck{}
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
	sortChecks(sentChecks, sortMethod, sortDir)
	c.JSON(200, gin.H{"checks": sentChecks, "count": len(CachedChecks.checks),
		"filteredCount": filteredCount, "updatedAt": CachedChecks.updatedAt, "types": CachedChecks.types})
}
