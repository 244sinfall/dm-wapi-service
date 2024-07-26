package logcleaner

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

func processLog(inputFile multipart.File) (*os.File, error) {
	scanner := bufio.NewScanner(inputFile)
	var removed int
	f, err := ioutil.TempFile("", "output-")
	if err != nil {
		return nil, errors.New("unable to create temp file on server")
	}
mainLoop:
	for scanner.Scan() {
		line := scanner.Text()

		for _, rule := range ignoreRules {
			if strings.Contains(line, rule) {
				removed++
				continue mainLoop
			}
		}
		_, err = f.WriteString(line + "\n")
		if err != nil {
			return nil, errors.New("unable to write cleaned log to temp file")
		}

	}
	_, err = f.WriteString(fmt.Sprintf("\n\nОчищено при помощи dm.rolevik.site. Удалено %v строк.", removed))
	if err != nil {
		return nil, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println("Error when closing temp file: ", err.Error(), "File: ", f.Name())
		}
	}()
	return f, nil
}
