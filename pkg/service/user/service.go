package user

import "github.com/Ayocodes24/GO-Eats/pkg/database"

type UsrService struct {
	db  database.Database
	env string
}

func NewUserService(db database.Database, env string) *UsrService {
	return &UsrService{db, env}
}
