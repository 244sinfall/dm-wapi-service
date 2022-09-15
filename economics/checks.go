package economics

import (
	"cloud.google.com/go/firestore"
	"context"
	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"gopkg.in/headzoo/surf.v1"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Check struct {
	Id          int      `firestore:"Id,omitempty"`
	Date        string   `firestore:"Date,omitempty"`
	Owner       string   `firestore:"Owner,omitempty"`
	CheckType   string   `firestore:"type,omitempty"`
	Money       int      `firestore:"Money,omitempty"`
	Name        string   `firestore:"Name,omitempty"`
	Description string   `firestore:"Description,omitempty"`
	Body        []string `firestore:"Body,omitempty"`
	Status      string   `firestore:"Status,omitempty"`
	Gm          string   `firestore:"Gm,omitempty"`
}

const defaultCheckCount = 13000

func ParseChecksFromDarkmoon() ([]Check, error) {
	browser := surf.NewBrowser()
	browser.SetUserAgent("PostmanRuntime/7.29.2")
	browser.SetTransport(cloudflarebp.AddCloudFlareByPass(&http.Transport{}))
	err := browser.Open("https://dm.rolevik.site/mock.html")
	if err != nil {
		return []Check{}, err
	}
	table := browser.State().Dom.Find("table") //Child 1 - headers, child 2 - contents
	checks := make([]Check, 0, defaultCheckCount)
	table.Children().Next().Children().Each(func(i int, sel *goquery.Selection) {
		data := sel.Children()
		container := data.Find(".container-items")
		if container.Size() != 0 {
			container.Children().Each(func(i int, item *goquery.Selection) {
				if len(checks) > 0 {
					checks[len(checks)-1].Body = append(checks[len(checks)-1].Body, item.Text())
				}
			})
		} else {
			checkCells := data.Contents().Nodes
			checkId, _ := strconv.Atoi(checkCells[0].Data)
			checkMoney, _ := strconv.Atoi(checkCells[4].Data)
			var gm, status, description string
			if len(checkCells) == 8 {
				status = checkCells[6].Data
				gm = checkCells[7].Data
			}
			if len(checkCells) == 9 {
				status = checkCells[7].Data
				gm = checkCells[8].Data
				description = checkCells[6].Data
			}
			checks = append(checks, Check{
				Id:          checkId,
				Date:        checkCells[1].Data,
				Owner:       checkCells[2].Data,
				CheckType:   checkCells[3].Data,
				Money:       checkMoney,
				Name:        checkCells[5].Data,
				Description: description,
				Body:        make([]string, 0, 1),
				Status:      status,
				Gm:          gm,
			})
		}
	})
	return checks, nil
}

func filterChecksStatus(c []Check, status string) []Check {
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
	newChecks := make([]Check, 0, len(c)/5)
	for _, check := range c {
		if check.Status == statusName {
			newChecks = append(newChecks, check)
		}
	}
	return newChecks
}

func sortChecks(c []Check, sortBy string, ascending bool) {
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

func filterChecks(c []Check, phrase string) []Check {
	if len(c) == 0 {
		return c
	}
	newChecks := make([]Check, 0, len(c)/10)
	for _, check := range c {
		if strings.Contains(strconv.Itoa(check.Id), phrase) || strings.Contains(check.Gm, phrase) ||
			strings.Contains(check.Name, phrase) || strings.Contains(check.Owner, phrase) ||
			strings.Contains(check.Description, phrase) {
			newChecks = append(newChecks, check)
		}
	}
	return newChecks
}

func ReceiveChecks(c *gin.Context, f *firestore.Client, ctx context.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip")) //
	search := c.Query("search")
	status := c.Query("status")
	sortMethod := c.Query("sortBy")           //
	sortDirection := c.Query("sortDirection") //
	force := c.Query("force")
	if CachedChecks.updating {
		c.JSON(500, gin.H{"error": "Checks are currently unavailable due to cache update"})
		return
	}
	if len(CachedChecks.checks) == 0 {
		ChecksScheduler(f, ctx, true)
	}
	if force != "" {
		if time.Now().Sub(CachedChecks.updatedAt) < 5*time.Minute {
			c.JSON(400, gin.H{"error": "Force update is available if cached checks are older than 5 minutes", "updatedAt": CachedChecks.updatedAt})
			return
		} else {
			err := ParseAndDeployNewChecks(f, ctx)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}
	sentChecks := CachedChecks.checks
	filteredCount := len(sentChecks)
	if status != "" {
		sentChecks = filterChecksStatus(sentChecks, status)
		filteredCount = len(sentChecks)
	}
	if search != "" {
		sentChecks = filterChecks(sentChecks, search)
		filteredCount = len(sentChecks)
	}
	if skip != 0 {
		if len(sentChecks)-1 > skip {
			sentChecks = sentChecks[skip:]
		} else {
			sentChecks = []Check{}
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
	c.JSON(200, gin.H{"checks": sentChecks, "count": len(CachedChecks.checks), "filteredCount": filteredCount, "updatedAt": CachedChecks.updatedAt})
}
