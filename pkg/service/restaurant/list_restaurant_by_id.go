package restaurant

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
)

func (restSrv *RestaurantService) ListRestaurantById(ctx context.Context, restaurantId int64) (restaurant.Restaurant, error) {
	var restro restaurant.Restaurant

	err := restSrv.db.Select(ctx, &restro, "restaurant_id", restaurantId)
	if err != nil {
		return restaurant.Restaurant{}, err
	}
	return restro, nil
}
