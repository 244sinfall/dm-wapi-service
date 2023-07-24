package users

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"net/smtp"
	"os"
)

const adminPermission = 4

func ResetUserPassword(c *gin.Context, a *auth.Client, f *firestore.Client, ctx context.Context) {
	token, err := a.VerifyIDToken(ctx, c.Request.Header.Get("Authorization"))
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
	link, err := a.PasswordResetLink(ctx, body.Email)
	from := "dm@244sinfall.ru"
	password := os.Getenv("DM_API_EMAIL_PASSWORD")
	to := []string{body.Email}
	host := "smtp.beget.com"
	port := "25"
	messageSubject := "Сброс пароля для Darkmoon WAPI\n"
	messageBody := "Для сброса пароля " + name + " перейдите по ссылке:\n" + link
	message := []byte(messageSubject + messageBody)
	emailAuth := smtp.PlainAuth("", from, password, host)
	err = smtp.SendMail(host+":"+port, emailAuth, from, to, message)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
	return
}
