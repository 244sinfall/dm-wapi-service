package common

import "strings"


var legitSuffixes = []string{
	" W",
	" D",
	" M",
	" WM",
	" MW",
	" WD",
	" DW",
}

func GetSuffixIfLegit(line string) (bool, string) {
	for _, suffix := range legitSuffixes {
		if strings.HasSuffix(line, suffix) {
			return true, suffix
		}
	}
	return false, ""
}