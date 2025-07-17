package delivery

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/delivery"
)

func (deliverSrv *DeliveryService) AddDeliveryPerson(ctx context.Context, deliveryPerson *delivery.DeliveryPerson) (bool, error) {
	_, err := deliverSrv.db.Insert(ctx, deliveryPerson)
	if err != nil {
		return false, err
	}
	return true, nil
}
