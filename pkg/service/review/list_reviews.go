package review

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/review"
)

func (revSrv *ReviewService) ListReviews(ctx context.Context, restaurantId int64) ([]review.Review, error) {
	var reviewList []review.Review

	err := revSrv.db.Select(ctx, &reviewList, "restaurant_id", restaurantId)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}
