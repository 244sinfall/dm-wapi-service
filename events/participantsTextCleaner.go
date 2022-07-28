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
		//fmt.Println(line)
		//empty check
		if len(strings.TrimSpace(line)) == 0 {
			//fmt.Println("empty line skipped")
			continue
		}
		lineHasLegitSuffix := checkForLegitSuffixes(line)
		if lineHasLegitSuffix {
			//fmt.Println("has suffix")
			if strings.Count(line, " ") == 1 {
				cleaned += line + "\n"
				//fmt.Println("has one spaced, added to " + cleaned)
				continue
			}
		}
		//fmt.Println("has no suffix")
		before, after, found := strings.Cut(line, " ")
		//fmt.Println("before: " + before + "after: " + after)
		regexp, _ := regexp2.Compile("([А-яА-Я])+")
		if found {
			//fmt.Println("reg exp work:")
			//fmt.Println(regexp.FindString(before))
			cleaned += regexp.FindString(before) + "\n"
			edited += before + " " + after + "\n"
		} else {
			//fmt.Println("reg exp work:")
			//fmt.Println(regexp.FindString(line))
			cleaned += regexp.FindString(line) + "\n"
		}
		continue

		//if strings.ContainsAny(line, " W") || strings.ContainsAny(line, " M") ||
		//	strings.ContainsAny(line, " W") ||  strings.ContainsAny(line, " W") ||
		//	strings.ContainsAny(line, " W")
	}
	c.JSON(http.StatusOK, gin.H{"cleanedText": cleaned, "editedLines": edited})
}
