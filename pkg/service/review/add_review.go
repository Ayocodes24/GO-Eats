package review

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/review"
)

func (revSrv *ReviewService) Add(ctx context.Context, review *review.Review) (bool, error) {
	_, err := revSrv.db.Insert(ctx, review)
	if err != nil {
		return false, err
	}
	return true, nil
}
