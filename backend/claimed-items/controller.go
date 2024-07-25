package claimeditems

import (
	"darkmoon-wapi-service/auth"
	"darkmoon-wapi-service/globals"
	rabbitmqclient "darkmoon-wapi-service/rabbitmq_client"
	rabbitmqmessages "darkmoon-wapi-service/rabbitmq_messages"

	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func AddClaimedItem(c *gin.Context) {
	user, err := auth.Authenticate(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	if !user.IsGM() {
		c.JSON(403, gin.H{"error": "You dont have permission"})
		return
	}
	claimedItem := new(claimedItem)
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(claimedItem)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = add(*claimedItem)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(200)
		baseMessage := rabbitmqmessages.GetAuthenticatedBaseMessage(user, "Claimed-Item", "Added")
		rabbitmqclient.SendLogMessage(GetClaimedItemRabbitMQMessage(baseMessage, claimedItem, nil), user, nil)
		return
	}
}

func ApproveClaimedItem(c *gin.Context) {
	user, err := auth.Authenticate(c.Request.Header.Get("Authorization"))
	a := globals.GetAuth()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !user.IsAdmin() {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	id := c.Param("id")
	fbUser, err := a.GetUser(globals.GetGlobalContext(), user.IntegrationUserId)
	if err != nil {
		c.JSON(500, gin.H{"error": "User not found"})
		return
	}
	err = approve(id, fbUser.DisplayName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		c.Status(200)
		return
	}
}

func UpdateClaimedItem(c *gin.Context) {
	user, err := auth.Authenticate(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !user.IsGM() {
		c.JSON(400, gin.H{"error": "You dont have permission"})
		return
	}
	fbUser, _ := globals.GetAuth().GetUser(globals.GetGlobalContext(), user.IntegrationUserId)
	claimedItemMock := new(claimedItem)
	decoder := json.NewDecoder(c.Request.Body)
	err = decoder.Decode(claimedItemMock)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	f := globals.GetFirestore()
	docRef := f.Doc("claimedItems/"+id)
	doc, err := docRef.Get(globals.GetGlobalContext())
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	var oldItem claimedItem
	err = doc.DataTo(&oldItem)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if !user.IsAdmin() && (oldItem.Name != claimedItemMock.Name || oldItem.Link != claimedItemMock.Link || oldItem.Reviewer != claimedItemMock.Reviewer) {
		c.JSON(403, gin.H{"error": "Not allowed"})
		return
	}
	if oldItem.Reviewer == claimedItemMock.Reviewer {
		y, m, d := time.Now().Date()
		claimedItemMock.Reviewer += fmt.Sprintf("\nИзменил: %v (%v.%v.%v)", fbUser.DisplayName, d, int(m), y)
	}
	claimedItemMock.Id = oldItem.Id
	
	newItem, err := update(id, *claimedItemMock)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		baseMessage := rabbitmqmessages.GetAuthenticatedBaseMessage(user, "Claimed-Item", "Updated")
		rabbitmqclient.SendLogMessage(GetClaimedItemRabbitMQMessage(baseMessage, newItem, &oldItem), user, nil)
		c.Status(200)
		return
	}
}

func DeleteClaimedItem(c *gin.Context) {
	user, err := auth.Authenticate(c.Request.Header.Get("Authorization"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if !user.IsAdmin() {
		c.JSON(403, gin.H{"error": "You dont have permission"})
		return
	}
	id := c.Param("id")
	item, err := delete(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	baseMessage := rabbitmqmessages.GetAuthenticatedBaseMessage(user, "Claimed-Item", "Deleted")
	rabbitmqclient.SendLogMessage(GetClaimedItemRabbitMQMessage(baseMessage, item, nil), user, nil)
	c.Status(200)
}

func ReceiveClaimedItems(c *gin.Context) {
	c.JSON(200, gin.H{"result": getClaimedItems()})
}
