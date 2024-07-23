package services

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type GameObject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

var gameObjects = make([]GameObject, 0, 120000)

func init() {
	f, err := os.Open(os.Getenv("GOBS_FILE_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range data {
		id, _ := strconv.Atoi(v[0])
		gobType, err := strconv.Atoi(v[2])
		if err != nil || gobType < 0 || gobType > 2 {
			continue
		}
		gameObjects = append(gameObjects, GameObject{
			Id:   id,
			Name: v[1],
			Type: gobType,
		})
	}
}

func GetGobs() []GameObject {
	return gameObjects
}