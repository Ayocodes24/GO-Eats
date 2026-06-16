package cart_order

import (
	"context"
	cartModel "github.com/Ayocodes24/GO-Eats/pkg/database/models/cart"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
)

func (cartSrv *CartService) ListItems(ctx context.Context, cartId int64) (*[]cartModel.CartItems, error) {
	var cartItems []cartModel.CartItems
	var relatedFields = []string{"Cart", "Restaurant", "MenuItem"}
	whereFilter := database.Filter{"cart_items.cart_id": cartId}
	err := cartSrv.db.SelectWithRelation(ctx, &cartItems, relatedFields, whereFilter)
	if err != nil {
		return nil, err
	}
	return &cartItems, nil
}
