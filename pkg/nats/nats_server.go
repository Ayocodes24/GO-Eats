package nats

import (
	"errors"
	"log"
	"log/slog"
	"strings"

	"github.com/Ayocodes24/GO-Eats/pkg/wsclients"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

type NATS struct {
	Conn *nats.Conn
}

var errNATSUnavailable = errors.New("NATS not connected")

func NewNATS(url string) (*NATS, error) {
	nc, err := nats.Connect(url, nats.Name("food-delivery-nats"))
	if err != nil {
		return nil, err
	}
	return &NATS{Conn: nc}, nil
}

func (n *NATS) Pub(topic string, message []byte) error {
	if n == nil || n.Conn == nil {
		return errNATSUnavailable
	}
	return n.Conn.Publish(topic, message)
}

func (n *NATS) Sub(topic string, clients *wsclients.Registry) error {
	if n == nil || n.Conn == nil {
		return errNATSUnavailable
	}
	_, err := n.Conn.Subscribe(topic, func(msg *nats.Msg) {
		message := string(msg.Data)
		slog.Info("MESSAGE_REPLY_FROM_NATS", "RECEIVED_MESSAGE", message)
		userId, messageData := n.formatMessage(message)
		if conn, ok := clients.Get(userId); ok {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(messageData)); err != nil {
				log.Println("Error sending message to client:", err)
				conn.Close()
				clients.Delete(userId)
			}
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (n *NATS) formatMessage(message string) (userId string, messageData string) {
	parts := strings.Split(message, "|")
	result := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2) // Split into key and value
		if len(kv) == 2 {
			result[kv[0]] = kv[1] // Store in a map
		}
	}
	return result["USER_ID"], result["MESSAGE"]
}
