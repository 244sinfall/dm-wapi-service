package gob


type GameObject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}