package events

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"net/http"
	regexp2 "regexp"
	"strings"
)

type Participants struct {
	RawText string `json:"RawText"`
}

var legitSuffixes = []string{
	" W",
	" D",
	" M",
	" WM",
	" MW",
	" WD",
	" DW",
}

func checkForLegitSuffixes(line string) bool {
	for _, suffix := range legitSuffixes {
		if strings.HasSuffix(line, suffix) {
			return true
		}
	}
	return false
}

func CleanParticipantsText(c *gin.Context) {
	var raw Participants
	var cleaned string
	var edited string
	err := c.BindJSON(&raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(raw.RawText))
	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		lineHasLegitSuffix := checkForLegitSuffixes(line)
		if lineHasLegitSuffix {
			if strings.Count(line, " ") == 1 {
				cleaned += line + "\n"
				continue
			}
		}
		before, after, found := strings.Cut(line, " ")
		regexp, _ := regexp2.Compile("([А-яА-Я])+")
		if found {
			cleaned += regexp.FindString(before) + "\n"
			edited += before + " " + after + "\n"
		} else {
			cleanedStr := regexp.FindString(line)
			cleanedLine := cleanedStr + "\n"
			cleaned += cleanedLine
			if cleanedStr != line {
				edited += line + "\n"
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"cleanedText": cleaned, "editedLines": edited})
}
