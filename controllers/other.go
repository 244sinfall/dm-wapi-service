package controllers

import (
	services "darkmoon-wapi-service/services"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxFileSize int64 = 1e+8

func LogClean(c *gin.Context) {
	file, err := c.FormFile("input")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !strings.HasSuffix(file.Filename, ".txt") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong file format!"})
		return
	}
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Max file size is 100mb"})
		return
	}
	textFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cleanedFile, err := services.CleanLog(textFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "utf8")
	c.Header("Content-Disposition", "attachment; filename="+cleanedFile.Name())
	c.Header("Content-Type", "application/octet-stream")
	c.FileAttachment(cleanedFile.Name(), "output.txt")
	defer func() {
		err := textFile.Close()
		if err != nil {
			fmt.Println("Error when closing downloaded file: ", err.Error())
		}
		err = os.Remove(cleanedFile.Name())
		if err != nil {
			fmt.Println("Error when deleting temp file: ", err.Error(), "File: ", cleanedFile.Name())
		}
	}()
}