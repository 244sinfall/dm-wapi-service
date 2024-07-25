package checks

import (
	"fmt"
)

type checkUser struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	GameId   int    `json:"gameId"`
}

type checkResponse struct {
	Types  []string             `json:"types"`
	Result []checkResponseCheck `json:"result"`
	Count  int                  `json:"count"`
}

type checkResponseCheck struct {
	Id       int       `json:"id"`
	Date     string    `json:"date"`
	Sender   checkUser `json:"senderUser"` // owner
	Receiver string    `json:"receiver"`   // checktype
	Subject  string    `json:"subject"`    // name
	Body     string    `json:"body"`       // description
	Money    int       `json:"money"`
	GmUser   checkUser `json:"gmUser"`
	Status   string    `json:"status"`
	Items    string    `json:"items"`
}

type checkRequestFilter struct {
	Limit    int
	Skip     int
	Search   string
	Category string
	Status   string
}

func (f *checkRequestFilter) ToCheckServiceQueryString() string {
	return fmt.Sprintf("?limit=%v&skip=%v&search=%s&type=%s&status=%s", f.Limit, f.Skip, f.Search, f.Category, f.Status)
}
