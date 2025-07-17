package announcements

import "github.com/Ayocodes24/GO-Eats/pkg/database"

type AnnouncementService struct {
	db  database.Database
	env string
}

func NewAnnouncementService(db database.Database, env string) *AnnouncementService {
	return &AnnouncementService{db, env}
}
