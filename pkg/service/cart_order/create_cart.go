package cart_order

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/abstract/cart"
)

func (cartSrv *CartService) Create(ctx context.Context, cart *cart.Cart) (*cart.Cart, error) {
	_, err := cartSrv.db.Insert(ctx, cart)
	if err != nil {
		return nil, err
	}
	return cart, nil
}
