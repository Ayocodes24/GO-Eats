package cart_order

import (
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/nats"
)

type CartService struct {
	db   database.Database
	env  string
	nats *nats.NATS
}

func NewCartService(db database.Database, env string, nats *nats.NATS) *CartService {
	return &CartService{db, env, nats}
}
