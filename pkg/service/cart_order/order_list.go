package cart_order

import (
	"context"
	"errors"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/order"
)

func (cartSrv *CartService) OrderList(ctx context.Context, userId int64) (*[]order.Order, error) {
	var ordersList []order.Order

	err := cartSrv.db.Select(ctx, &ordersList, "user_id", userId)
	if err != nil {
		return nil, err
	}
	return &ordersList, nil
}

func (cartSrv *CartService) OrderItemsList(ctx context.Context, userId int64, orderId int64) (*[]order.OrderItems, error) {
	var ordersItemsList []order.OrderItems

	count, err := cartSrv.db.Count(ctx, "orders", "COUNT(*)", "user_id", userId)
	if err != nil {
		return nil, errors.New("invalid order")
	}

	if count == 0 {
		return nil, errors.New("invalid order. No order found")
	}

	var relatedFields = []string{"Restaurant", "MenuItem"}
	whereFilter := database.Filter{"order_id": orderId}

	err = cartSrv.db.SelectWithRelation(ctx, &ordersItemsList, relatedFields, whereFilter)

	if err != nil {
		return nil, err
	}

	return &ordersItemsList, nil
}
