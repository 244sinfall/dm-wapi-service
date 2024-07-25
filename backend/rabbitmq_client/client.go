package rabbitmqclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"darkmoon-wapi-service/auth"
	rabbitmqmessages "darkmoon-wapi-service/rabbitmq_messages"

	amqp "github.com/rabbitmq/amqp091-go"
)

var rabbitMqConnection *amqp.Connection = nil

func init() {
	conn, err := amqp.DialTLS(os.Getenv("DM_API_RABBITMQ_STRING"), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatalf("Cannot connect to rabbitmq: %v\n", err)
	}
	rabbitMqConnection = conn
}

func SendLogMessage(message rabbitmqmessages.IRabbitMQMessage, user *auth.WapiAuthenticatedUser, target *auth.WapiAuthenticatedUser) {
	ch, err := rabbitMqConnection.Channel()
	if err != nil {
		fmt.Printf("Cannot create a channel: %v\n", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",  // name
		"topic", // type
		false,   // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		fmt.Printf("Cannot declare an exchange: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Cannot marshal message: %v\n", err)
		return
	}
	err = ch.PublishWithContext(ctx,
		"logs", // exchange
		message.GetTopic(),
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		fmt.Printf("Cannot publish message: %v\n", err)
		return
	}

	log.Printf(" [x] Sent %s\n", body)
}
