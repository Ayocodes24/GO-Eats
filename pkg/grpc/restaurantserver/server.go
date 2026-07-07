package restaurantserver

import (
	"context"
	"fmt"

	restaurantpb "github.com/Ayocodes24/GO-Eats/pkg/proto/restaurant"
	restaurantsvc "github.com/Ayocodes24/GO-Eats/pkg/service/restaurant"
)

type Server struct {
	restaurantpb.UnimplementedRestaurantServiceServer
	svc *restaurantsvc.RestaurantService
}

func New(svc *restaurantsvc.RestaurantService) *Server {
	return &Server{svc: svc}
}

// GetMenuItem fetches a single menu item's pricing details.
// Called by Order Service when computing order totals at checkout.
func (s *Server) GetMenuItem(ctx context.Context, req *restaurantpb.GetMenuItemRequest) (*restaurantpb.GetMenuItemResponse, error) {
	menus, err := s.svc.ListAllMenus(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list menus: %w", err)
	}
	for _, m := range menus {
		if m.MenuID == req.MenuId {
			return &restaurantpb.GetMenuItemResponse{
				MenuId:       m.MenuID,
				Name:         m.Name,
				Price:        m.Price,
				RestaurantId: m.RestaurantID,
			}, nil
		}
	}
	return nil, fmt.Errorf("menu item %d not found", req.MenuId)
}
