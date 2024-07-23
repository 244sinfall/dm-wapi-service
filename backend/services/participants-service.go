package services

import (
	"bufio"
	regexp2 "regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ParticipantsRequest struct {
	RawText string `json:"rawText"  binding:"required"`
}

type ParticipantsResponse struct {
	CleanedText string `json:"cleanedText"`
	Count       int    `json:"cleanedCount"`
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

func (t ParticipantsRequest) CleanRawText() ParticipantsResponse {
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
		regexp, _ := regexp2.Compile("([А-яА-ЯЁё])+( [WDM]{1,2})?")
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
	r.Count = count
	return r
}
