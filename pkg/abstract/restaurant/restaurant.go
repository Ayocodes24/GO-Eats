package restaurant

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
)

type Restaurant interface {
	Add(ctx context.Context, user *restaurant.Restaurant) (bool, error)
	ListRestaurants(ctx context.Context) ([]restaurant.Restaurant, error)
	ListRestaurantById(ctx context.Context, restaurantId int64) (restaurant.Restaurant, error)
	DeleteRestaurant(ctx context.Context, restaurantId int64) (bool, error)
}
