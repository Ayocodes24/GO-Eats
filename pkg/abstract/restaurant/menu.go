package restaurant

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/restaurant"
)

type MenuItems interface {
	AddMenu(ctx context.Context, menu *restaurant.MenuItem) (*restaurant.MenuItem, error)
	UpdateMenuPhoto(ctx context.Context, menu *restaurant.MenuItem)
	ListMenus(ctx context.Context, restaurantId int64) ([]restaurant.MenuItem, error)
	ListAllMenus(ctx context.Context) ([]restaurant.MenuItem, error)
	DeleteMenu(ctx context.Context, menuId int64, restaurantId int64) (bool, error)
}
