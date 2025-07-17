package cart_order

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
)

func (cartSrv *CartService) DeleteItem(ctx context.Context, cartItemId int64) (bool, error) {
	filter := database.Filter{"cart_item_id": cartItemId}

	_, err := cartSrv.db.Delete(ctx, "cart_items", filter)
	if err != nil {
		return false, err
	}
	return true, nil
}
