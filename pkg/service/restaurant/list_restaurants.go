package restaurant

import (
	"context"
	restaurantModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
)

func (restSrv *RestaurantService) ListRestaurants(ctx context.Context) ([]restaurantModel.Restaurant, error) {
	var restaurants []restaurantModel.Restaurant

	err := restSrv.db.SelectAll(ctx, "restaurant", &restaurants)
	if err != nil {
		return nil, err
	}

	return restaurants, nil
}
