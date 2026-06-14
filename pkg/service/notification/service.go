package notification

import (
	"log/slog"

	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/nats"
	"github.com/Ayocodes24/GO-Eats/pkg/wsclients"
)

type NotificationService struct {
	db   database.Database
	env  string
	nats *nats.NATS
}

func NewNotificationService(db database.Database, env string, nats *nats.NATS) *NotificationService {
	return &NotificationService{db, env, nats}
}

func (s *NotificationService) SubscribeNewOrders(clients *wsclients.Registry) error {
	slog.Info("Listening to ==> NotificationService::SubscribeNewOrders")
	return s.nats.Sub("orders.new.*", clients)
}

func (s *NotificationService) SubscribeOrderStatus(clients *wsclients.Registry) error {
	slog.Info("Listening to ==> NotificationService::SubscribeOrderStatus")
	return s.nats.Sub("orders.status.*", clients)
}
