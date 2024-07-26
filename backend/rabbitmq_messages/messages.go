package rabbitmqmessages

import (
	"darkmoon-wapi-service/auth"
	"fmt"
	"strings"
)

type IRabbitMQMessage interface {
	GetTopic() string
}

type RabbitMQMessage struct {
	Service   string `json:"service"`
	Entity    string `json:"entity"`
	Operation string `json:"operation"`
}

func (m *RabbitMQMessage) GetTopic() string {
	return fmt.Sprintf("logs.wapi-service.%s.%s", strings.ToLower(m.Entity), strings.ToLower(m.Operation))
}

func GetBaseMessage(entity string, operation string) *RabbitMQMessage {
	var new = new(RabbitMQMessage)
	new.Service = "wapi-service"
	new.Entity = entity
	new.Operation = operation
	return new
}

type AuthenticatedRabbitMQMessage struct {
	*RabbitMQMessage
	Id                int    `json:"id"`
	Integration       string `json:"integration"`
	IntegrationUserId string `json:"integrationUserId"`
}

func (m *AuthenticatedRabbitMQMessage) GetTopic() string {
	return fmt.Sprintf("logs.web-authorized.wapi-service.%s.%s", strings.ToLower(m.Entity), strings.ToLower(m.Operation))
}

func GetAuthenticatedBaseMessage(au *auth.WapiAuthenticatedUser, entity string, operation string) *AuthenticatedRabbitMQMessage {
	var new = new(AuthenticatedRabbitMQMessage)
	new.RabbitMQMessage = GetBaseMessage(entity, operation)
	new.Id = au.UserId
	new.Integration = "Wapi"
	new.IntegrationUserId = au.IntegrationUserId
	return new
}
