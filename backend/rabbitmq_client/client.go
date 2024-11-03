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

	"github.com/getsentry/sentry-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMqConnection struct {
	conn  *amqp.Connection
	retry bool
}

// var rabbitMqConnection *amqp.Connection = nil
var client *rabbitMqConnection = nil

func init() {
	client = new(rabbitMqConnection)
	client.createConnection()
	client.allowRetry()
}

func (c *rabbitMqConnection) allowRetry() {
	c.retry = true
}

func (c *rabbitMqConnection) isAbleToRetry() bool {
	return c.retry
}

func (c *rabbitMqConnection) restrictRetry() {
	c.retry = false
}

func (c *rabbitMqConnection) createConnection() {
	conn, err := amqp.DialTLS(os.Getenv("DM_API_RABBITMQ_STRING"), &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatalf("Cannot connect to rabbitmq: %v\n", err)
	}
	c.conn = conn
}

func (c *rabbitMqConnection) getChannel() *amqp.Channel {
	channel, err := c.conn.Channel()
	if err != nil {
		if c.isAbleToRetry() {
			sentry.CaptureMessage("Error on getting a channel: " + err.Error() + " Retrying!")
			c.createConnection()
			c.restrictRetry()
			return c.getChannel()
		} else {
			sentry.CaptureException(err)
			return nil
		}
	}
	c.allowRetry()
	return channel
}

func SendLogMessage(message rabbitmqmessages.IRabbitMQMessage, user *auth.WapiAuthenticatedUser, target *auth.WapiAuthenticatedUser) {
	ch := client.getChannel()
	defer ch.Close()

	err := ch.ExchangeDeclare(
		"logs",  // name
		"topic", // type
		false,   // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		sentry.CaptureException(err)
		fmt.Printf("Cannot declare an exchange: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(message)
	if err != nil {
		sentry.CaptureException(err)
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
		sentry.CaptureException(err)
		fmt.Printf("Cannot publish message: %v\n", err)
		return
	}

	log.Printf(" [x] Sent %s\n", body)
}
