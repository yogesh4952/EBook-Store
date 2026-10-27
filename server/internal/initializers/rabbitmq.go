package initializers

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yogesh4952/ebookstore/pkg/utils"
)

type RabbitMqPayload struct {
	Ch   *amqp.Channel
	Conn *amqp.Connection
}

func InitRabbitmq() RabbitMqPayload {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}

	conn, err := amqp.Dial(url)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"orders", // name
		"topic",  // type
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // no-wait
		nil,      // args
	)
	utils.FailOnError(err, "Failed to declare exchange")

	return RabbitMqPayload{Ch: ch, Conn: conn}
}
