package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CheckUser struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	GameId   int    `json:"gameId"`
}

type CheckResponse struct {
	Types  []string             `json:"types"`
	Result []CheckResponseCheck `json:"result"`
	Count  int                  `json:"count"`
}

type CheckResponseCheckItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type CheckResponseCheck struct {
	Id       int       `json:"id"`
	Date     string    `json:"date"`
	Sender   CheckUser `json:"senderUser"` // owner
	Receiver string    `json:"receiver"`   // checktype
	Subject  string    `json:"subject"`    // name
	Body     string    `json:"body"`       // description
	Money    int       `json:"money"`
	GmUser   CheckUser `json:"gmUser"`
	Status   string    `json:"status"`
	Items    string    `json:"items"`
}

const defaultCheckCount = 13000
const cacheFrequency = 30 * time.Minute

func parseChecksFromDarkmoon() (*CheckResponse, error) {
	newChecksResponse := new(CheckResponse)

	newChecksResponse.Result = make([]CheckResponseCheck, 0, defaultCheckCount)
	response, err := http.Get(os.Getenv("DM_API_CHECKS_ADDRESS"))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&newChecksResponse)
	if len(newChecksResponse.Result) == 0 {
		return nil, errors.New("nothing parsed")
	}
	if err != nil {
		return nil, err
	}
	return newChecksResponse, nil
}

func FilterChecksCategory(c []CheckResponseCheck, category string) []CheckResponseCheck {
	newChecks := make([]CheckResponseCheck, 0, len(c)/5)
	for _, check := range c {
		if check.Receiver == category {
			newChecks = append(newChecks, check)
		}
	}
	return newChecks
}

func FilterChecksStatus(c []CheckResponseCheck, status string) []CheckResponseCheck {
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
	newChecks := make([]CheckResponseCheck, 0, len(c)/5)
	for _, check := range c {
		if check.Status == statusName {
			newChecks = append(newChecks, check)
		}
	}
	return newChecks
}

func SortChecks(c []CheckResponseCheck, sortBy string, ascending bool) {
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

func FilterChecks(c []CheckResponseCheck, phrase string) []CheckResponseCheck {
	if len(c) == 0 {
		return c
	}
	newChecks := make([]CheckResponseCheck, 0, len(c)/10)
	lowerPhrase := strings.ToLower(phrase)
	for _, check := range c {
		if strings.Contains(strconv.Itoa(check.Id), phrase) || findTextMatch(check.GmUser.Nickname, lowerPhrase) ||
			findTextMatch(check.Subject, lowerPhrase) || findTextMatch(check.Sender.Nickname, lowerPhrase) ||
			findTextMatch(check.Body, lowerPhrase) || findTextMatch(check.Receiver, lowerPhrase) {
			newChecks = append(newChecks, check)
		}
	}
	return newChecks
}

type CachedChecks struct {
	Checks    []CheckResponseCheck `json:"checks"`
	UpdatedAt time.Time            `json:"updatedAt"`
	Types     []string             `json:"types"`
	Updating  bool
}

var cachedChecks CachedChecks

func GetCachedChecks() CachedChecks {
	return cachedChecks
}

func findCheckTypes(c []CheckResponseCheck) []string {
	checkTypes := make(map[string]struct{}, 10)
	for _, v := range c {
		checkTypes[v.Receiver] = struct{}{}
	}
	checkTypesSlice := make([]string, 0, len(checkTypes))
	for k := range checkTypes {
		checkTypesSlice = append(checkTypesSlice, k)
	}
	sort.Strings(checkTypesSlice)
	return checkTypesSlice
}

func ParseAndDeployNewChecks() error {
	cachedChecks.Updating = true
	parsedChecks, err := parseChecksFromDarkmoon()
	if err != nil {
		fmt.Println("Unable to parse checks from Darkmoon. Error: " + err.Error())
		cachedChecks.Updating = false
		return err
	} else {
		cachedChecks = CachedChecks{
			Checks:    parsedChecks.Result,
			UpdatedAt: time.Now(),
			Types:     parsedChecks.Types,
			Updating:  false,
		}
	}
	return nil
}

func ChecksScheduler(ping bool) {
	fmt.Println("Scheduler just started. Retrieving checks")
	for {
		if time.Now().Sub(cachedChecks.UpdatedAt) > cacheFrequency {
			err := ParseAndDeployNewChecks()
			if err != nil {
				fmt.Println("Unable to parse new checks! " + err.Error())
			}
			err = nil
			time.Sleep(5 * time.Minute)
		} else {
			if ping {
				return
			}
			schedule := cacheFrequency - time.Now().Sub(cachedChecks.UpdatedAt) + (2 * time.Minute)
			time.Sleep(schedule)
		}
	}
}
