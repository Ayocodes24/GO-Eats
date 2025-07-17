package review

import "github.com/Ayocodes24/GO-Eats/pkg/database"

type ReviewService struct {
	db  database.Database
	env string
}

func NewReviewService(db database.Database, env string) *ReviewService {
	return &ReviewService{db, env}
}
