package review

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
)

func (revSrv *ReviewService) DeleteReview(ctx context.Context, reviewId int64, userId int64) (bool, error) {
	filter := database.Filter{"review_id": reviewId, "user_id": userId}

	_, err := revSrv.db.Delete(ctx, "reviews", filter)
	if err != nil {
		return false, err
	}
	return true, nil
}
