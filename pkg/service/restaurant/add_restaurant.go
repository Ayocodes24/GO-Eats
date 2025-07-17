package restaurant

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
)

func (restSrv *RestaurantService) Add(ctx context.Context, restaurant *restaurant.Restaurant) (bool, error) {
	_, err := restSrv.db.Insert(ctx, restaurant)
	if err != nil {
		return false, err
	}
	return true, nil
}
