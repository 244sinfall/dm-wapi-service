package events

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	regexp2 "regexp"
	"strings"
)

type ParticipantsRequest struct {
	RawText string `json:"rawText"  binding:"required"`
}

type ParticipantsResponse struct {
	CleanedText string `json:"cleanedText"`
	EditedLines string `json:"editedLines"`
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

func CheckForLegitSuffixes(line string) (bool, string) {
	for _, suffix := range legitSuffixes {
		if strings.HasSuffix(line, suffix) {
			return true, suffix
		}
	}
	return false, ""
}

func cleanRawText(t ParticipantsRequest) ParticipantsResponse {
	var r ParticipantsResponse
	var count int
	scanner := bufio.NewScanner(strings.NewReader(t.RawText))
	for scanner.Scan() {
		line := scanner.Text()
		// Empty line handler
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		// Capitalize name + regexp handler
		regexp, _ := regexp2.Compile("([А-яА-Я])+( [WDM]{1,2})?")
		foundRegexp := regexp.FindString(line)
		makeTitle := cases.Title(language.Russian)
		name, suffix, found := strings.Cut(foundRegexp, " ")
		foundRegexp = makeTitle.String(name)
		if found {
			foundRegexp += " " + suffix
		}
		if strings.Contains(r.CleanedText, name) {
			r.EditedLines += foundRegexp + " (дубликат)\n"
			continue
		}
		if foundRegexp != line {
			r.EditedLines += line + " (изменено под формат)\n"
		}
		if !found || (suffix != "M" && suffix != "WM" && suffix != "MW") {
			count++
		}
		r.CleanedText += foundRegexp + "\n"
	}
	if r.EditedLines != "" {
		r.EditedLines = "Обработка соответствия строк формату:\n" + r.EditedLines + "\n"
	}
	if count < 5 {
		var newCleanedText string
		var newEditedLines string
		scanner = bufio.NewScanner(strings.NewReader(r.CleanedText))
		for scanner.Scan() {
			participant := scanner.Text()
			_, suffix := CheckForLegitSuffixes(participant)
			if suffix == " M" || suffix == " MW" || suffix == " WM" {
				newEditedLines += participant + " (недостаточно участников)\n"
				continue
			}
			if suffix == " DW" || suffix == " WD" {
				newCleanedText += strings.TrimSuffix(participant, suffix) + " D\n"
				newEditedLines += participant + " (недостаточно участников для бонуса)\n"
				continue
			}
			if suffix == " W" {
				newCleanedText += strings.TrimSuffix(participant, suffix) + "\n"
				newEditedLines += participant + " (недостаточно участников для бонуса)\n"
				continue
			}
			newCleanedText += participant + "\n"
		}
		r.CleanedText = newCleanedText
		if newEditedLines != "" {
			r.EditedLines += "Обработка соблюдения условий:\n" + newEditedLines
		}
	}
	return r
}

func CleanParticipantsText(c *gin.Context) {
	var raw ParticipantsRequest
	err := c.BindJSON(&raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	response := cleanRawText(raw)
	c.JSON(http.StatusOK, response)
}
