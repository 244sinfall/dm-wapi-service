package checks

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

func getChecks(filter *checkRequestFilter) (*checkResponse, error) {
	newChecksResponse := new(checkResponse)

	newChecksResponse.Result = make([]checkResponseCheck, 0, filter.Limit)
	response, err := http.Get(os.Getenv("DM_API_CHECKS_ADDRESS")+filter.ToCheckServiceQueryString())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
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