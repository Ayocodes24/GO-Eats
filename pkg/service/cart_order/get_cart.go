package cart_order

import (
	"context"
	cartModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/cart"
)

func (cartSrv *CartService) GetCartId(ctx context.Context, UserId int64) (*cartModel.Cart, error) {
	var cartInfo cartModel.Cart

	err := cartSrv.db.Select(ctx, &cartInfo, "user_id", UserId)
	if err != nil {
		return nil, err
	}
	return &cartInfo, nil
}
