package cart_order

import (
	"context"
	cartModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/cart"
)

func (cartSrv *CartService) Create(ctx context.Context, cart *cartModel.Cart) (*cartModel.Cart, error) {
	_, err := cartSrv.db.Insert(ctx, cart)
	if err != nil {
		return nil, err
	}
	return cart, nil
}
