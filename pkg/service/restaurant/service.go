package restaurant

import "github.com/Ayocodes24/GO-Eats/pkg/database"

type RestaurantService struct {
	db  database.Database
	env string
}

func NewRestaurantService(db database.Database, env string) *RestaurantService {
	return &RestaurantService{db, env}
}
