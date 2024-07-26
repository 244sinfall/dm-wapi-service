package checks

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

var checkHttpClient = &http.Client{}

func getChecks(filter *checkRequestFilter) (*checkResponse, error) {
	newChecksResponse := new(checkResponse)
	newChecksResponse.Result = make([]checkResponseCheck, 0, filter.Limit)
	newChecksResponse.Types = make([]string, 0, 10)
	url := os.Getenv("DM_API_CHECKS_ADDRESS") + filter.ToCheckServiceQueryString()
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Content-Type", "application/json")
	request.Close = true
	if err != nil {
		return nil, err
	}
	fmt.Println(url)
	fmt.Println(request)
	response, err := checkHttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&newChecksResponse)
	response.Body.Close()
	if err != nil {
		fmt.Println("Error decoding checks: " + err.Error())
		return nil, err
	}
	if len(newChecksResponse.Result) == 0 {
		return nil, errors.New("nothing parsed")
	}
	return newChecksResponse, nil
}
