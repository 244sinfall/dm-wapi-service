package rabbitmqmessages

import services "darkmoon-wapi-service/services"

type ClaimedItemRabbitMQMessage struct {
	*AuthenticatedRabbitMQMessage
	Item    *services.ClaimedItem
	OldItem *services.ClaimedItem
}

func GetClaimedItemRabbitMQMessage(base *AuthenticatedRabbitMQMessage, item *services.ClaimedItem, oldItem *services.ClaimedItem) *ClaimedItemRabbitMQMessage {
	new := &ClaimedItemRabbitMQMessage{AuthenticatedRabbitMQMessage: base, Item: item, OldItem: oldItem}
	return new
}
