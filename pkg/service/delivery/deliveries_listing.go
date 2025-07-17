package delivery

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/delivery"
)

func (deliverSrv *DeliveryService) DeliveryListing(ctx context.Context, orderID int64, userID int64) (*[]delivery.Deliveries, error) {
	var deliveriesList []delivery.Deliveries
	whereFilter := database.Filter{"order_id": orderID, "delivery_person_id": userID}
	err := deliverSrv.db.SelectWithMultipleFilter(ctx, &deliveriesList, whereFilter)
	if err != nil {
		return nil, err
	}
	return &deliveriesList, nil
}
