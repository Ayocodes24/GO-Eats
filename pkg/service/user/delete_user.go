package user

import (
	"context"
	"github.com/Ayocodes24/GO-Eats/pkg/database"
)

func (usrSrv *UsrService) Delete(ctx context.Context, userId int64) (bool, error) {
	filter := database.Filter{"id": userId}
	_, err := usrSrv.db.Delete(ctx, "users", filter)
	if err != nil {
		return false, err
	}
	return true, nil
}
