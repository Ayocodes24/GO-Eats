package restaurant

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
)

func (restSrv *RestaurantService) DeleteRestaurant(ctx context.Context, restaurantId int64) (bool, error) {
	filter := database.Filter{"restaurant_id": restaurantId}

	_, err := restSrv.db.Delete(ctx, "restaurant", filter)
	if err != nil {
		return false, err
	}
	return true, nil
}
