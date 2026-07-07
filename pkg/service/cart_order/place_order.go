package cart_order

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/cart"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/order"
)

func (cartSrv *CartService) PlaceOrder(ctx context.Context, cartId int64, userId int64, deliveryAddress string) (*order.Order, error) {
	var cartItems []cart.CartItems
	var newOrder order.Order
	var newOrderItems []order.OrderItems
	var orderTotal float64

	whereFilter := database.Filter{"cart_items.cart_id": cartId}

	if cartSrv.priceFetcher != nil {
		// Microservice mode: fetch cart items without cross-DB relation.
		if err := cartSrv.db.SelectWithMultipleFilter(ctx, &cartItems, whereFilter); err != nil {
			return nil, err
		}
	} else {
		// Monolith mode: resolve MenuItem via DB join.
		if err := cartSrv.db.SelectWithRelation(ctx, &cartItems, []string{"MenuItem"}, whereFilter); err != nil {
			return nil, err
		}
	}

	if len(cartItems) == 0 {
		return nil, errors.New("no items in cart")
	}

	newOrder.UserID = userId
	newOrder.OrderStatus = "pending"
	newOrder.TotalAmount = 0
	newOrder.DeliveryAddress = deliveryAddress

	if _, err := cartSrv.db.Insert(ctx, &newOrder); err != nil {
		return nil, err
	}

	newOrderItems = make([]order.OrderItems, len(cartItems))
	for i, cartItem := range cartItems {
		var price float64
		if cartSrv.priceFetcher != nil {
			p, err := cartSrv.priceFetcher.GetPrice(ctx, cartItem.ItemID)
			if err != nil {
				return nil, fmt.Errorf("could not fetch price for item %d: %w", cartItem.ItemID, err)
			}
			price = p * float64(cartItem.Quantity)
		} else {
			price = cartItem.MenuItem.Price * float64(cartItem.Quantity)
		}

		newOrderItems[i].OrderID = newOrder.OrderID
		newOrderItems[i].ItemID = cartItem.ItemID
		newOrderItems[i].RestaurantID = cartItem.RestaurantID
		newOrderItems[i].Quantity = cartItem.Quantity
		newOrderItems[i].Price = price

		if _, err := cartSrv.db.Insert(ctx, &newOrderItems[i]); err != nil {
			return nil, err
		}
		orderTotal += price
	}

	if _, err := cartSrv.db.Update(ctx, "orders",
		database.Filter{"total_amount": orderTotal, "order_status": "in_progress"},
		database.Filter{"order_id": newOrder.OrderID},
	); err != nil {
		return nil, err
	}

	return &newOrder, nil
}

func (cartSrv *CartService) RemoveItemsFromCart(ctx context.Context, cartId int64) error {
	filter := database.Filter{"cart_id": cartId}
	if _, err := cartSrv.db.Delete(ctx, "cart_items", filter); err != nil {
		return errors.New("failed to delete cart items")
	}
	return nil
}

func (cartSrv *CartService) NewOrderPlacedNotification(userId int64, orderId int64) error {
	message := fmt.Sprintf("USER_ID:%d|MESSAGE:Your order number %d has been successfully placed, and the chef has begun the cooking process.", userId, orderId)
	topic := fmt.Sprintf("orders.new.%d", userId)
	return cartSrv.nats.Pub(topic, []byte(message))
}
