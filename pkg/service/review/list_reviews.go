package review

func (revSrv *ReviewService) ListReviews(ctx context.Context, restaurantId int64) ([]review.Review, error) {
	var reviewList []review.Review

	err := revSrv.db.Select(ctx, &reviewList, "restaurant_id", restaurantId)
	if err != nil {
		return nil, err
	}
	return reviewList, nil
}
