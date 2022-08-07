package events

import (
	"bufio"
	"github.com/gin-gonic/gin"
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
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if strings.Contains(r.CleanedText, line) {
			r.EditedLines += line + " (дубликат)\n"
			continue
		}
		lineHasLegitSuffix, suffix := CheckForLegitSuffixes(line)
		if lineHasLegitSuffix && strings.Count(line, " ") == 1 {
			if strings.Contains(r.CleanedText, strings.TrimSuffix(line, suffix)) {
				r.EditedLines += line + " (дубликат)\n"
				continue
			}
			if suffix != " M" && suffix != " WM" && suffix != " MW" {
				count++
			}
			r.CleanedText += line + "\n"
			continue
		}
		before, after, found := strings.Cut(line, " ")
		regexp, _ := regexp2.Compile("([А-яА-Я])+")
		if found {
			fixed := regexp.FindString(before)
			if strings.Contains(r.CleanedText, fixed) {
				r.EditedLines += fixed + " (дубликат)\n"
				continue
			}
			count++
			r.CleanedText += fixed + "\n"
			r.EditedLines += before + " " + after + "\n"
		} else {
			cleanedStr := regexp.FindString(line)
			if strings.Contains(r.CleanedText, cleanedStr) {
				r.EditedLines += cleanedStr + " (дубликат)\n"
				continue
			}
			cleanedLine := cleanedStr + "\n"
			r.CleanedText += cleanedLine
			count++
			if cleanedStr != line {
				r.EditedLines += line + "\n"
			}
		}
	}
	if count < 5 {
		var newCleanedText string
		scanner = bufio.NewScanner(strings.NewReader(r.CleanedText))
		for scanner.Scan() {
			participant := scanner.Text()
			if strings.Contains(newCleanedText, participant) {
				r.EditedLines += participant + " (дубликат)\n"
				continue
			}
			_, suffix := CheckForLegitSuffixes(participant)
			if suffix == " M" || suffix == " MW" || suffix == " WM" {
				r.EditedLines += participant + " (недостаточно участников)\n"
				continue
			}
			if suffix == " DW" || suffix == " WD" {
				newCleanedText += strings.TrimSuffix(participant, suffix) + " D\n"
				r.EditedLines += participant + " (недостаточно участников для бонуса)\n"
				continue
			}
			if suffix == " W" {
				newCleanedText += strings.TrimSuffix(participant, suffix) + "\n"
				r.EditedLines += participant + " (недостаточно участников для бонуса)"
				continue
			}
			newCleanedText += participant + "\n"
		}
		r.CleanedText = newCleanedText
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
