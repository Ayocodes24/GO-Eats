package delivery

import "github.com/Ayocodes24/GO-Eats/pkg/database"

type DeliveryService struct {
	db   database.Database
	env  string
	nats *nats.NATS
}

func NewDeliveryService(db database.Database, env string, nats *nats.NATS) *DeliveryService {
	return &DeliveryService{db, env, nats}
}
