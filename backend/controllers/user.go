package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/mail.v2"
)

func ResetUserPassword(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	permInfo, err := f.Doc("permissions/" + token.UID).Get(ctx)
	permData := permInfo.Data()
	permission := permData["permission"].(int64)
	name := permData["name"].(string)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if permission < adminPermission {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	body := new(struct {
		Email string `json:"email"`
	})
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(body)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	link, err := a.PasswordResetLink(ctx, body.Email)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	from := "dm@244sinfall.ru"
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", body.Email)
	m.SetHeader("Subject", "Сброс пароля для Darkmoon WAPI")
	m.SetHeader("Message-ID", fmt.Sprintf("%s@%s", uuid.New().String(), "244sinfall.ru"))
	m.SetBody("text/plain", "Для сброса пароля "+name+" перейдите по ссылке:\n"+link)
	password := os.Getenv("DM_API_EMAIL_PASSWORD")
	host := "smtp.beget.com"
	port := 25
	d := mail.NewDialer(host, port, from, password)
	err = d.DialAndSend(m)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}
