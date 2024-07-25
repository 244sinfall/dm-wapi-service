package claimeditems

import (
	rabbitmqmessages "darkmoon-wapi-service/rabbitmq_messages"
	"time"
)

type claimedItem struct {
	Id             string    `json:"id"`
	Quality        string    `json:"quality"`
	Name           string    `json:"name"`
	Link           string    `json:"link"`
	Owner          string    `json:"owner"`
	OwnerProfile   string    `json:"ownerProfile"`
	OwnerProofName string    `json:"ownerProof"`
	OwnerProofLink string    `json:"ownerProofLink"`
	Reviewer       string    `json:"reviewer"`
	Accepted       bool      `json:"accepted"`
	Acceptor       string    `json:"acceptor"`
	AddedAt        time.Time `json:"addedAt"`
	AcceptedAt     time.Time `json:"acceptedAt"`
	AdditionalInfo string    `json:"additionalInfo"`
}

func (c *claimedItem) GetKey() string {
	switch c.Quality {
	case "Легендарный":
		return "legendary"
	case "Эпический":
		return "epic"
	case "Редкий":
		return "rare"
	case "Необычный":
		return "green"
	default:
		return "other"
	}
}

type ClaimedItemRabbitMQMessage struct {
	*rabbitmqmessages.AuthenticatedRabbitMQMessage
	Item    *claimedItem
	OldItem *claimedItem
}

func GetClaimedItemRabbitMQMessage(base *rabbitmqmessages.AuthenticatedRabbitMQMessage, item *claimedItem, oldItem *claimedItem) *ClaimedItemRabbitMQMessage {
	new := &ClaimedItemRabbitMQMessage{AuthenticatedRabbitMQMessage: base, Item: item, OldItem: oldItem}
	return new
}
