package gob

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/csv"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
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

const gmPermission = 1

func init() {
	f, err := os.Open("gobs.csv")
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
		id, err := strconv.Atoi(v[0])
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
func ReceiveGobs(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, _ := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	permInfo, _ := f.Doc("permissions/" + token.UID).Get(ctx)
	permission := permInfo.Data()["permission"].(int64)
	if permission < gmPermission {
		c.JSON(400, gin.H{"error": "Not enough permissions"})
		return
	}
	c.JSON(200, gin.H{"result": gameObjects})
}
