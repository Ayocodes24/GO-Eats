package review

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database/models/review"
)

type Review interface {
	Add(ctx context.Context, review *review.Review) (bool, error)
	ListReviews(ctx context.Context, restaurantId int64) ([]review.Review, error)
	DeleteReview(ctx context.Context, reviewId int64, userId int64) (bool, error)
}
