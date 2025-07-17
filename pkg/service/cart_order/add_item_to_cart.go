package cart_order

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/abstract/cart"
)

func (cartSrv *CartService) AddItem(ctx context.Context, Item *cart.CartItems) (*cart.CartItems, error) {
	_, err := cartSrv.db.Insert(ctx, Item)
	if err != nil {
		return nil, err
	}
	return Item, nil
}
