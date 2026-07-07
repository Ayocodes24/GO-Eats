package cart_order

import (
	"context"

	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/nats"
)

// MenuPriceFetcher lets the service fetch a menu item's price without
// coupling to a specific database or gRPC client.
// In the monolith, this is backed by a direct DB query.
// In the Order microservice, this is backed by a gRPC call to Restaurant Service.
type MenuPriceFetcher interface {
	GetPrice(ctx context.Context, menuID int64) (float64, error)
}

type CartService struct {
	db           database.Database
	env          string
	nats         *nats.NATS
	priceFetcher MenuPriceFetcher
}

func NewCartService(db database.Database, env string, nats *nats.NATS) *CartService {
	return &CartService{db: db, env: env, nats: nats}
}

func NewCartServiceWithFetcher(db database.Database, env string, nats *nats.NATS, pf MenuPriceFetcher) *CartService {
	return &CartService{db: db, env: env, nats: nats, priceFetcher: pf}
}
