package main

import (
	"context"
	"fmt"

	restaurantpb "github.com/Ayocodes24/GO-Eats/pkg/proto/restaurant"
)

// grpcPriceFetcher implements cart_order.MenuPriceFetcher using the
// Restaurant Service gRPC client. Order Service calls this instead of
// querying the restaurant database directly.
type grpcPriceFetcher struct {
	client restaurantpb.RestaurantServiceClient
}

func (f *grpcPriceFetcher) GetPrice(ctx context.Context, menuID int64) (float64, error) {
	resp, err := f.client.GetMenuItem(ctx, &restaurantpb.GetMenuItemRequest{MenuId: menuID})
	if err != nil {
		return 0, fmt.Errorf("restaurant-svc gRPC GetMenuItem(%d): %w", menuID, err)
	}
	return resp.Price, nil
}
